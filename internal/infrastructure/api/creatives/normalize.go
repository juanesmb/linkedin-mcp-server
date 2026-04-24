package creatives

import (
	"fmt"
	"strconv"
	"strings"
)

const sponsoredCreativeURNPrefix = "urn:li:sponsoredCreative:"

// NormalizeCreative maps one LinkedIn creative element to the stable MCP DTO.
func NormalizeCreative(raw map[string]any) NormalizedCreative {
	out := NormalizedCreative{}

	if raw == nil {
		return out
	}

	out.CreativeURN, out.CreativeID = normalizeCreativeID(raw["id"])
	out.CampaignURN = stringFromAny(raw["campaign"])
	out.IntendedStatus = stringFromAny(raw["intendedStatus"])
	out.Format = stringFromAny(raw["format"])

	if serving, ok := raw["isServing"].(bool); ok {
		out.IsServing = &serving
	}

	if review, ok := raw["review"].(map[string]any); ok {
		out.ReviewStatus = stringFromAny(review["status"])
	}

	content, _ := raw["content"].(map[string]any)
	inline, _ := raw["inlineContent"].(map[string]any)

	switch {
	case content != nil:
		out.ContentKind, out.Headline, out.Description, out.CTA, out.LandingPageURL = extractFromContentMap(content)
	case inline != nil:
		out.ContentKind, out.Headline, out.Description, out.CTA, out.LandingPageURL = extractFromInlineContent(inline)
	default:
		out.ContentKind = "unknown"
	}

	return out
}

func normalizeCreativeID(id any) (urn string, numericID string) {
	switch v := id.(type) {
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			return "", ""
		}
		if strings.HasPrefix(s, sponsoredCreativeURNPrefix) {
			return s, strings.TrimPrefix(s, sponsoredCreativeURNPrefix)
		}
		return sponsoredCreativeURNPrefix + s, s
	case float64:
		n := int64(v)
		num := strconv.FormatInt(n, 10)
		return sponsoredCreativeURNPrefix + num, num
	case int64:
		num := strconv.FormatInt(v, 10)
		return sponsoredCreativeURNPrefix + num, num
	default:
		return "", ""
	}
}

func stringFromAny(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case float64:
		return strconv.FormatInt(int64(t), 10)
	case bool:
		return strconv.FormatBool(t)
	default:
		return strings.TrimSpace(fmt.Sprint(t))
	}
}

func extractFromContentMap(content map[string]any) (kind, headline, description, cta, landing string) {
	if ref, ok := content["reference"].(string); ok && strings.TrimSpace(ref) != "" {
		return "content_reference", "", "", "", ""
	}
	if m, ok := content["textAd"].(map[string]any); ok {
		return "text_ad",
			stringFromAny(m["headline"]),
			stringFromAny(m["description"]),
			"",
			stringFromAny(m["landingPage"])
	}
	if m, ok := content["spotlight"].(map[string]any); ok {
		return "spotlight",
			stringFromAny(m["headline"]),
			stringFromAny(m["description"]),
			stringFromAny(m["callToAction"]),
			stringFromAny(m["landingPage"])
	}
	if m, ok := content["followCompany"].(map[string]any); ok {
		return "follow_company",
			stringFromAny(m["headline"]),
			stringFromAny(m["description"]),
			stringFromAny(m["callToAction"]),
			""
	}
	if m, ok := content["jobPosting"].(map[string]any); ok {
		return "job_posting",
			stringFromAny(m["title"]),
			stringFromAny(m["description"]),
			"",
			stringFromAny(m["landingPage"])
	}
	if m, ok := content["article"].(map[string]any); ok {
		return "article",
			stringFromAny(m["title"]),
			stringFromAny(m["description"]),
			"",
			stringFromAny(m["landingPage"])
	}
	if m, ok := content["document"].(map[string]any); ok {
		return "document",
			stringFromAny(m["title"]),
			"",
			"",
			stringFromAny(m["landingPage"])
	}
	// Sponsored update / media content often uses nested shapes; try common keys.
	for _, key := range []string{"sponsoredUpdate", "sponsoredContent", "media"} {
		if m, ok := content[key].(map[string]any); ok {
			h, d, c, l := extractHeadlineDescriptionCTALanding(m)
			if h != "" || d != "" || l != "" || c != "" {
				return key, h, d, c, l
			}
		}
	}
	return "structured_content", "", "", "", ""
}

func extractFromInlineContent(inline map[string]any) (kind, headline, description, cta, landing string) {
	if post, ok := inline["post"].(map[string]any); ok {
		kind = "inline_post"
		headline = firstNonEmpty(
			stringFromAny(post["title"]),
			stringFromAny(post["subject"]),
			stringFromAny(post["commentary"]),
		)
		description = stringFromAny(post["commentary"])
		if description == headline {
			description = ""
		}
		cta = stringFromAny(post["contentCallToActionLabel"])
		landing = stringFromAny(post["contentLandingPage"])
		return kind, headline, description, cta, landing
	}
	return "inline_other", "", "", "", ""
}

func extractHeadlineDescriptionCTALanding(m map[string]any) (headline, description, cta, landing string) {
	headline = firstNonEmpty(
		stringFromAny(m["headline"]),
		stringFromAny(m["title"]),
		stringFromAny(m["subject"]),
	)
	description = stringFromAny(m["description"])
	cta = stringFromAny(m["callToAction"])
	landing = firstNonEmpty(
		stringFromAny(m["landingPage"]),
		stringFromAny(m["landingPageUrl"]),
		stringFromAny(m["contentLandingPage"]),
	)
	if media, ok := m["media"].(map[string]any); ok {
		if headline == "" {
			headline = stringFromAny(media["title"])
		}
	}
	return headline, description, cta, landing
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
