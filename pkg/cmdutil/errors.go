// Copyright (c) 2026 CrowdStrike, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmdutil

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-openapi/runtime"
)

// fqlSyntaxRe matches a minimal valid FQL clause: field + operator + value.
// A bare word like "severity" with no operator is invalid FQL.
var fqlSyntaxRe = regexp.MustCompile(`[:<>!~]`)

// ValidateFQLFilter returns an error if the filter string looks like malformed FQL
// (e.g. a bare field name with no operator or value).
func ValidateFQLFilter(filter string) error {
	if filter == "" {
		return nil
	}
	if !fqlSyntaxRe.MatchString(filter) {
		return fmt.Errorf("invalid FQL filter %q — missing operator and value\n"+
			"  FQL requires field + operator + value, e.g.:\n"+
			"    --filter \"severity:'high'\"\n"+
			"    --filter \"severity:'high'+status:'open'\"", filter)
	}
	return nil
}

// scopeRequirements maps gofalcon operation names to the Falcon API scopes they require.
// Scope names match exactly what is shown in the Falcon console API client editor
// (https://falcon.crowdstrike.com/api-clients-and-keys).
//
// Note: ClusterCombined and ContainerCombined use the KubernetesProtection gofalcon
// client but actually call /container-security/* endpoints — they require
// "Falcon Container Image: Read", NOT "Kubernetes Protection: Read".
var scopeRequirements = map[string]string{
	// Assets
	"CloudSecurityAssetsQueries":     "CSPM Assessment: Read",
	"CloudSecurityAssetsEntitiesGet": "CSPM Assessment: Read",
	// Risks
	"CombinedCloudRisks": "CSPM Assessment: Read",
	// IOM
	"CspmEvaluationsIomQueries":  "CSPM Assessment: Read",
	"CspmEvaluationsIomEntities": "CSPM Assessment: Read",
	// IaC
	"CombinedDetections": "IaC Security: Read",
	// Compliance
	"QueryComplianceFrameworks":                "CSPM Registration: Read",
	"QueryComplianceControls":                  "CSPM Registration: Read",
	"CloudComplianceFrameworkPostureSummaries": "CSPM Registration: Read",
	"CloudComplianceRulePostureSummaries":      "CSPM Registration: Read",
	// Suppression
	"QuerySuppressionRules":  "CSPM Registration: Read",
	"GetSuppressionRules":    "CSPM Registration: Read",
	"CreateSuppressionRule":  "CSPM Registration: Write",
	"DeleteSuppressionRules": "CSPM Registration: Write",
	// Groups — /cloud-security/entities/cloud-groups/v1 requires asset-groups:write
	"CreateCloudGroupExternal":    "Asset Groups: Write",
	"UpdateCloudGroupExternal":    "Asset Groups: Write",
	"DeleteCloudGroupsExternal":   "Asset Groups: Write",
	"ListCloudGroupsByIDExternal": "Asset Groups: Read",
	"ListCloudGroupsExternal":     "Asset Groups: Read",
	// Kubernetes clusters + containers — /container-security/* endpoints
	"ClusterCombined":   "Falcon Container Image: Read",
	"ContainerCombined": "Falcon Container Image: Read",
	// Container images
	"CombinedImageDetail": "Falcon Container Image: Read",
	// Vulnerabilities
	"ReadCombinedVulnerabilities": "Falcon Container Image: Read",

	// Human-readable operation strings passed to HandleAPIError by commands.
	"query cloud assets":                         "CSPM Assessment: Read",
	"get cloud assets":                           "CSPM Assessment: Read",
	"list cloud risks":                           "CSPM Assessment: Read",
	"query IOM findings":                         "CSPM Assessment: Read",
	"get IOM findings":                           "CSPM Assessment: Read",
	"query compliance frameworks":                "CSPM Registration: Read",
	"query compliance controls":                  "CSPM Registration: Read",
	"get compliance framework posture summaries": "CSPM Registration: Read",
	"get compliance rule posture summaries":      "CSPM Registration: Read",
	"query suppression rules":                    "CSPM Registration: Read",
	"get suppression rules":                      "CSPM Registration: Read",
	"update suppression rule":                    "CSPM Registration: Write",
}

// cwppOperations are served by cwppcontainersecurityapi. A 403 on these
// may mean either a missing scope OR the CWPP product not being provisioned.
var cwppOperations = map[string]bool{
	"ClusterCombined":             true,
	"ContainerCombined":           true,
	"CombinedImageDetail":         true,
	"ReadCombinedVulnerabilities": true,
}

// featureRequirements maps operations to features that must be tenant-provisioned.
// Keys are the human-readable operation strings passed to HandleAPIError.
var featureRequirements = map[string]string{
	// Compliance posture — both the ID-discovery query and the posture summary fetch
	// return 404 when the Compliance posture feature is not provisioned.
	"query compliance frameworks":                "CSPM Compliance posture feature (contact your Falcon admin to enable)",
	"get compliance framework posture summaries": "CSPM Compliance posture feature (contact your Falcon admin to enable)",
	"query compliance controls":                  "CSPM Compliance posture feature (contact your Falcon admin to enable)",
	"get compliance rule posture summaries":      "CSPM Compliance posture feature (contact your Falcon admin to enable)",
	// Legacy gofalcon operation-name keys kept for backward compatibility.
	"CloudComplianceFrameworkPostureSummaries": "CSPM Compliance posture feature (contact your Falcon admin to enable)",
	"CloudComplianceRulePostureSummaries":      "CSPM Compliance posture feature (contact your Falcon admin to enable)",
}

// filterHints provides known filter field names for APIs whose 400 responses
// do not include a machine-readable field list. Keyed by the operation string
// passed to HandleAPIError.
var filterHints = map[string]string{
	"ContainerCombined": "cluster_name, namespace, cloud, container_id, " +
		"privileged, running_as_root, allow_privilege_escalation, image_name",
}

// HandleAPIError wraps a gofalcon API error with an actionable message.
// It inspects the HTTP status code via runtime.APIError and maps errors
// to the specific API scope required or missing tenant features.
func HandleAPIError(err error, operation string) error {
	if err == nil {
		return nil
	}

	var apiErr *runtime.APIError
	if !errors.As(err, &apiErr) {
		// Not a typed API error — return as-is with operation context
		return fmt.Errorf("%s: %w", operation, err)
	}

	switch {
	case apiErr.Code == 401:
		return fmt.Errorf("%s: authentication failed (HTTP 401)\n"+
			"  Check that FALCON_CLIENT_ID and FALCON_CLIENT_SECRET are correct.\n"+
			"  For non-US-1 regions, set FALCON_BASE_URL (e.g. https://api.eu-1.crowdstrike.com).\n"+
			"  Run 'falcon --help' to see all credential options", operation)

	case apiErr.Code == 403:
		if scope, ok := scopeRequirements[operation]; ok {
			hint := fmt.Sprintf("%s: access denied (HTTP 403)\n"+
				"  Required scope: %s\n\n"+
				"  To add this scope to your API client:\n"+
				"    1. Open https://falcon.crowdstrike.com/api-clients-and-keys\n"+
				"    2. Select your API client → Edit\n"+
				"    3. Enable: %s\n"+
				"    4. Save — update FALCON_CLIENT_SECRET with the new secret",
				operation, scope, scope)
			if cwppOperations[operation] {
				hint += "\n\n  Note: If the scope is already enabled, this command also requires\n" +
					"  the Falcon Cloud Workload Protection (CWPP) product to be provisioned\n" +
					"  in your tenant. Contact your Falcon admin if unsure."
			}
			return fmt.Errorf("%s", hint)
		}
		return fmt.Errorf("%s: access denied (HTTP 403)\n"+
			"  Your API client is missing a required scope.\n"+
			"  Check scopes at: https://falcon.crowdstrike.com/api-clients-and-keys",
			operation)

	case apiErr.Code == 404:
		if feature, ok := featureRequirements[operation]; ok {
			return fmt.Errorf("%s: feature not available (HTTP 404)\n  Requires: %s", operation, feature)
		}
		return fmt.Errorf("%s: resource not found (HTTP 404)", operation)

	case apiErr.Code == 400:
		if hint, ok := filterHints[operation]; ok && strings.Contains(err.Error(), "{}") {
			return fmt.Errorf("%s: invalid filter expression\n  Known filter fields: %s", operation, hint)
		}
		return fmt.Errorf("%s: bad request (HTTP 400) — %s", operation, FirstLine(err.Error()))

	case apiErr.Code == 429:
		return fmt.Errorf("%s: rate limited (HTTP 429) — too many requests, wait and retry", operation)

	case apiErr.Code >= 500:
		return fmt.Errorf("%s: Falcon API server error — try again or check https://status.crowdstrike.com", operation)
	}

	return fmt.Errorf("%s: %w", operation, err)
}

func contains(s string, keywords ...string) bool {
	lower := strings.ToLower(s)
	for _, k := range keywords {
		if strings.Contains(lower, strings.ToLower(k)) {
			return true
		}
	}
	return false
}

// ContainsAnyFold reports whether s contains any of the keywords, case-insensitively.
func ContainsAnyFold(s string, keywords ...string) bool { return contains(s, keywords...) }

// FirstLine returns the first non-empty line of a string, trimmed.
func FirstLine(s string) string {
	for line := range strings.SplitSeq(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			// Trim gofalcon's verbose error prefix like "[POST /path][403] opName  &{..."
			if idx := strings.Index(line, "&{"); idx > 0 {
				line = strings.TrimSpace(line[:idx])
			}
			return line
		}
	}
	return s
}
