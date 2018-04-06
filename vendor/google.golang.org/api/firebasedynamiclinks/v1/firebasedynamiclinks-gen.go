// Package firebasedynamiclinks provides access to the Firebase Dynamic Links API.
//
// See https://firebase.google.com/docs/dynamic-links/
//
// Usage example:
//
//   import "google.golang.org/api/firebasedynamiclinks/v1"
//   ...
//   firebasedynamiclinksService, err := firebasedynamiclinks.New(oauthHttpClient)
package firebasedynamiclinks // import "google.golang.org/api/firebasedynamiclinks/v1"

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	context "golang.org/x/net/context"
	ctxhttp "golang.org/x/net/context/ctxhttp"
	gensupport "google.golang.org/api/gensupport"
	googleapi "google.golang.org/api/googleapi"
)

// Always reference these packages, just in case the auto-generated code
// below doesn't.
var _ = bytes.NewBuffer
var _ = strconv.Itoa
var _ = fmt.Sprintf
var _ = json.NewDecoder
var _ = io.Copy
var _ = url.Parse
var _ = gensupport.MarshalJSON
var _ = googleapi.Version
var _ = errors.New
var _ = strings.Replace
var _ = context.Canceled
var _ = ctxhttp.Do

const apiId = "firebasedynamiclinks:v1"
const apiName = "firebasedynamiclinks"
const apiVersion = "v1"
const basePath = "https://firebasedynamiclinks-ipv6.googleapis.com/"

// OAuth2 scopes used by this API.
const (
	// View and administer all your Firebase data and settings
	FirebaseScope = "https://www.googleapis.com/auth/firebase"
)

func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &Service{client: client, BasePath: basePath}
	s.ManagedShortLinks = NewManagedShortLinksService(s)
	s.ShortLinks = NewShortLinksService(s)
	s.V1 = NewV1Service(s)
	return s, nil
}

type Service struct {
	client    *http.Client
	BasePath  string // API endpoint base URL
	UserAgent string // optional additional User-Agent fragment

	ManagedShortLinks *ManagedShortLinksService

	ShortLinks *ShortLinksService

	V1 *V1Service
}

func (s *Service) userAgent() string {
	if s.UserAgent == "" {
		return googleapi.UserAgent
	}
	return googleapi.UserAgent + " " + s.UserAgent
}

func NewManagedShortLinksService(s *Service) *ManagedShortLinksService {
	rs := &ManagedShortLinksService{s: s}
	return rs
}

type ManagedShortLinksService struct {
	s *Service
}

func NewShortLinksService(s *Service) *ShortLinksService {
	rs := &ShortLinksService{s: s}
	return rs
}

type ShortLinksService struct {
	s *Service
}

func NewV1Service(s *Service) *V1Service {
	rs := &V1Service{s: s}
	return rs
}

type V1Service struct {
	s *Service
}

// AnalyticsInfo: Tracking parameters supported by Dynamic Link.
type AnalyticsInfo struct {
	// GooglePlayAnalytics: Google Play Campaign Measurements.
	GooglePlayAnalytics *GooglePlayAnalytics `json:"googlePlayAnalytics,omitempty"`

	// ItunesConnectAnalytics: iTunes Connect App Analytics.
	ItunesConnectAnalytics *ITunesConnectAnalytics `json:"itunesConnectAnalytics,omitempty"`

	// ForceSendFields is a list of field names (e.g. "GooglePlayAnalytics")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "GooglePlayAnalytics") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *AnalyticsInfo) MarshalJSON() ([]byte, error) {
	type NoMethod AnalyticsInfo
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// AndroidInfo: Android related attributes to the Dynamic Link.
type AndroidInfo struct {
	// AndroidFallbackLink: Link to open on Android if the app is not
	// installed.
	AndroidFallbackLink string `json:"androidFallbackLink,omitempty"`

	// AndroidLink: If specified, this overrides the ‘link’ parameter on
	// Android.
	AndroidLink string `json:"androidLink,omitempty"`

	// AndroidMinPackageVersionCode: Minimum version code for the Android
	// app. If the installed app’s version
	// code is lower, then the user is taken to the Play Store.
	AndroidMinPackageVersionCode string `json:"androidMinPackageVersionCode,omitempty"`

	// AndroidPackageName: Android package name of the app.
	AndroidPackageName string `json:"androidPackageName,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AndroidFallbackLink")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AndroidFallbackLink") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *AndroidInfo) MarshalJSON() ([]byte, error) {
	type NoMethod AndroidInfo
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CreateManagedShortLinkRequest: Request to create a managed Short
// Dynamic Link.
type CreateManagedShortLinkRequest struct {
	// DynamicLinkInfo: Information about the Dynamic Link to be
	// shortened.
	// [Learn
	// more](https://firebase.google.com/docs/reference/dynamic-links/link-sh
	// ortener).
	DynamicLinkInfo *DynamicLinkInfo `json:"dynamicLinkInfo,omitempty"`

	// LongDynamicLink: Full long Dynamic Link URL with desired query
	// parameters specified.
	// For
	// example,
	// "https://sample.app.goo.gl/?link=http://www.google.com&apn=co
	// m.sample",
	// [Learn
	// more](https://firebase.google.com/docs/reference/dynamic-links/link-sh
	// ortener).
	LongDynamicLink string `json:"longDynamicLink,omitempty"`

	// Name: Link name to associate with the link. It's used for marketer to
	// identify
	// manually-created links in the Firebase
	// console
	// (https://console.firebase.google.com/).
	// Links must be named to be tracked.
	Name string `json:"name,omitempty"`

	// Suffix: Short Dynamic Link suffix. Optional.
	Suffix *Suffix `json:"suffix,omitempty"`

	// ForceSendFields is a list of field names (e.g. "DynamicLinkInfo") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DynamicLinkInfo") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *CreateManagedShortLinkRequest) MarshalJSON() ([]byte, error) {
	type NoMethod CreateManagedShortLinkRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CreateManagedShortLinkResponse: Response to create a short Dynamic
// Link.
type CreateManagedShortLinkResponse struct {
	// ManagedShortLink: Short Dynamic Link value. e.g.
	// https://abcd.app.goo.gl/wxyz
	ManagedShortLink *ManagedShortLink `json:"managedShortLink,omitempty"`

	// PreviewLink: Preview link to show the link flow chart. (debug info.)
	PreviewLink string `json:"previewLink,omitempty"`

	// Warning: Information about potential warnings on link creation.
	Warning []*DynamicLinkWarning `json:"warning,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "ManagedShortLink") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ManagedShortLink") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *CreateManagedShortLinkResponse) MarshalJSON() ([]byte, error) {
	type NoMethod CreateManagedShortLinkResponse
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CreateShortDynamicLinkRequest: Request to create a short Dynamic
// Link.
type CreateShortDynamicLinkRequest struct {
	// DynamicLinkInfo: Information about the Dynamic Link to be
	// shortened.
	// [Learn
	// more](https://firebase.google.com/docs/reference/dynamic-links/link-sh
	// ortener).
	DynamicLinkInfo *DynamicLinkInfo `json:"dynamicLinkInfo,omitempty"`

	// LongDynamicLink: Full long Dynamic Link URL with desired query
	// parameters specified.
	// For
	// example,
	// "https://sample.app.goo.gl/?link=http://www.google.com&apn=co
	// m.sample",
	// [Learn
	// more](https://firebase.google.com/docs/reference/dynamic-links/link-sh
	// ortener).
	LongDynamicLink string `json:"longDynamicLink,omitempty"`

	// Suffix: Short Dynamic Link suffix. Optional.
	Suffix *Suffix `json:"suffix,omitempty"`

	// ForceSendFields is a list of field names (e.g. "DynamicLinkInfo") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DynamicLinkInfo") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *CreateShortDynamicLinkRequest) MarshalJSON() ([]byte, error) {
	type NoMethod CreateShortDynamicLinkRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CreateShortDynamicLinkResponse: Response to create a short Dynamic
// Link.
type CreateShortDynamicLinkResponse struct {
	// PreviewLink: Preview link to show the link flow chart. (debug info.)
	PreviewLink string `json:"previewLink,omitempty"`

	// ShortLink: Short Dynamic Link value. e.g.
	// https://abcd.app.goo.gl/wxyz
	ShortLink string `json:"shortLink,omitempty"`

	// Warning: Information about potential warnings on link creation.
	Warning []*DynamicLinkWarning `json:"warning,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "PreviewLink") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "PreviewLink") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *CreateShortDynamicLinkResponse) MarshalJSON() ([]byte, error) {
	type NoMethod CreateShortDynamicLinkResponse
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// DesktopInfo: Desktop related attributes to the Dynamic Link.
type DesktopInfo struct {
	// DesktopFallbackLink: Link to open on desktop.
	DesktopFallbackLink string `json:"desktopFallbackLink,omitempty"`

	// ForceSendFields is a list of field names (e.g. "DesktopFallbackLink")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DesktopFallbackLink") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *DesktopInfo) MarshalJSON() ([]byte, error) {
	type NoMethod DesktopInfo
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// DeviceInfo: Signals associated with the device making the request.
type DeviceInfo struct {
	// DeviceModelName: Device model name.
	DeviceModelName string `json:"deviceModelName,omitempty"`

	// LanguageCode: Device language code setting.
	LanguageCode string `json:"languageCode,omitempty"`

	// LanguageCodeFromWebview: Device language code setting obtained by
	// executing JavaScript code in
	// WebView.
	LanguageCodeFromWebview string `json:"languageCodeFromWebview,omitempty"`

	// LanguageCodeRaw: Device language code raw setting.
	// iOS does returns language code in different format than iOS
	// WebView.
	// For example WebView returns en_US, but iOS returns en-US.
	// Field below will return raw value returned by iOS.
	LanguageCodeRaw string `json:"languageCodeRaw,omitempty"`

	// ScreenResolutionHeight: Device display resolution height.
	ScreenResolutionHeight int64 `json:"screenResolutionHeight,omitempty,string"`

	// ScreenResolutionWidth: Device display resolution width.
	ScreenResolutionWidth int64 `json:"screenResolutionWidth,omitempty,string"`

	// Timezone: Device timezone setting.
	Timezone string `json:"timezone,omitempty"`

	// ForceSendFields is a list of field names (e.g. "DeviceModelName") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DeviceModelName") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *DeviceInfo) MarshalJSON() ([]byte, error) {
	type NoMethod DeviceInfo
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// DynamicLinkEventStat: Dynamic Link event stat.
type DynamicLinkEventStat struct {
	// Count: The number of times this event occurred.
	Count int64 `json:"count,omitempty,string"`

	// Event: Link event.
	//
	// Possible values:
	//   "DYNAMIC_LINK_EVENT_UNSPECIFIED" - Unspecified type.
	//   "CLICK" - Indicates that an FDL is clicked by users.
	//   "REDIRECT" - Indicates that an FDL redirects users to fallback
	// link.
	//   "APP_INSTALL" - Indicates that an FDL triggers an app install from
	// Play store, currently
	// it's impossible to get stats from App store.
	//   "APP_FIRST_OPEN" - Indicates that the app is opened for the first
	// time after an install
	// triggered by FDLs
	//   "APP_RE_OPEN" - Indicates that the app is opened via an FDL for
	// non-first time.
	Event string `json:"event,omitempty"`

	// Platform: Requested platform.
	//
	// Possible values:
	//   "DYNAMIC_LINK_PLATFORM_UNSPECIFIED" - Unspecified platform.
	//   "ANDROID" - Represents Android platform.
	// All apps and browsers on Android are classfied in this category.
	//   "IOS" - Represents iOS platform.
	// All apps and browsers on iOS are classfied in this category.
	//   "DESKTOP" - Represents desktop.
	// Note: other platforms like Windows, Blackberry, Amazon fall into
	// this
	// category.
	Platform string `json:"platform,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Count") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Count") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *DynamicLinkEventStat) MarshalJSON() ([]byte, error) {
	type NoMethod DynamicLinkEventStat
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// DynamicLinkInfo: Information about a Dynamic Link.
type DynamicLinkInfo struct {
	// AnalyticsInfo: Parameters used for tracking. See all tracking
	// parameters in
	// the
	// [documentation](https://firebase.google.com/docs/dynamic-links/cre
	// ate-manually).
	AnalyticsInfo *AnalyticsInfo `json:"analyticsInfo,omitempty"`

	// AndroidInfo: Android related information. See Android related
	// parameters in
	// the
	// [documentation](https://firebase.google.com/docs/dynamic-links/cre
	// ate-manually).
	AndroidInfo *AndroidInfo `json:"androidInfo,omitempty"`

	// DesktopInfo: Desktop related information. See desktop related
	// parameters in
	// the
	// [documentation](https://firebase.google.com/docs/dynamic-links/cre
	// ate-manually).
	DesktopInfo *DesktopInfo `json:"desktopInfo,omitempty"`

	// DomainUriPrefix: E.g. https://maps.app.goo.gl,
	// https://maps.page.link, https://g.co/maps
	// More examples can be found in description of getNormalizedUriPrefix
	// in
	// j/c/g/firebase/dynamiclinks/uri/DdlDomain.java
	DomainUriPrefix string `json:"domainUriPrefix,omitempty"`

	// DynamicLinkDomain: Dynamic Links domain that the project owns, e.g.
	// abcd.app.goo.gl
	// [Learn
	// more](https://firebase.google.com/docs/dynamic-links/android/receive)
	//
	// on how to set up Dynamic Link domain associated with your Firebase
	// project.
	//
	// Required.
	DynamicLinkDomain string `json:"dynamicLinkDomain,omitempty"`

	// IosInfo: iOS related information. See iOS related parameters in
	// the
	// [documentation](https://firebase.google.com/docs/dynamic-links/cre
	// ate-manually).
	IosInfo *IosInfo `json:"iosInfo,omitempty"`

	// Link: The link your app will open, You can specify any URL your app
	// can handle.
	// This link must be a well-formatted URL, be properly URL-encoded, and
	// use
	// the HTTP or HTTPS scheme. See 'link' parameters in
	// the
	// [documentation](https://firebase.google.com/docs/dynamic-links/cre
	// ate-manually).
	//
	// Required.
	Link string `json:"link,omitempty"`

	// NavigationInfo: Information of navigation behavior of a Firebase
	// Dynamic Links.
	NavigationInfo *NavigationInfo `json:"navigationInfo,omitempty"`

	// SocialMetaTagInfo: Parameters for social meta tag params.
	// Used to set meta tag data for link previews on social sites.
	SocialMetaTagInfo *SocialMetaTagInfo `json:"socialMetaTagInfo,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AnalyticsInfo") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AnalyticsInfo") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *DynamicLinkInfo) MarshalJSON() ([]byte, error) {
	type NoMethod DynamicLinkInfo
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// DynamicLinkStats: Analytics stats of a Dynamic Link for a given
// timeframe.
type DynamicLinkStats struct {
	// LinkEventStats: Dynamic Link event stats.
	LinkEventStats []*DynamicLinkEventStat `json:"linkEventStats,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "LinkEventStats") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "LinkEventStats") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *DynamicLinkStats) MarshalJSON() ([]byte, error) {
	type NoMethod DynamicLinkStats
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// DynamicLinkWarning: Dynamic Links warning messages.
type DynamicLinkWarning struct {
	// WarningCode: The warning code.
	//
	// Possible values:
	//   "CODE_UNSPECIFIED" - Unknown code.
	//   "NOT_IN_PROJECT_ANDROID_PACKAGE_NAME" - The Android package does
	// not match any in developer's DevConsole project.
	//   "NOT_INTEGER_ANDROID_PACKAGE_MIN_VERSION" - The Android minimum
	// version code has to be a valid integer.
	//   "UNNECESSARY_ANDROID_PACKAGE_MIN_VERSION" - Android package min
	// version param is not needed, e.g. when
	// 'apn' is missing.
	//   "NOT_URI_ANDROID_LINK" - Android link is not a valid URI.
	//   "UNNECESSARY_ANDROID_LINK" - Android link param is not needed, e.g.
	// when param 'al' and 'link' have
	// the same value..
	//   "NOT_URI_ANDROID_FALLBACK_LINK" - Android fallback link is not a
	// valid URI.
	//   "BAD_URI_SCHEME_ANDROID_FALLBACK_LINK" - Android fallback link has
	// an invalid (non http/https) URI scheme.
	//   "NOT_IN_PROJECT_IOS_BUNDLE_ID" - The iOS bundle ID does not match
	// any in developer's DevConsole project.
	//   "NOT_IN_PROJECT_IPAD_BUNDLE_ID" - The iPad bundle ID does not match
	// any in developer's DevConsole project.
	//   "UNNECESSARY_IOS_URL_SCHEME" - iOS URL scheme is not needed, e.g.
	// when 'ibi' are 'ipbi' are all missing.
	//   "NOT_NUMERIC_IOS_APP_STORE_ID" - iOS app store ID format is
	// incorrect, e.g. not numeric.
	//   "UNNECESSARY_IOS_APP_STORE_ID" - iOS app store ID is not needed.
	//   "NOT_URI_IOS_FALLBACK_LINK" - iOS fallback link is not a valid URI.
	//   "BAD_URI_SCHEME_IOS_FALLBACK_LINK" - iOS fallback link has an
	// invalid (non http/https) URI scheme.
	//   "NOT_URI_IPAD_FALLBACK_LINK" - iPad fallback link is not a valid
	// URI.
	//   "BAD_URI_SCHEME_IPAD_FALLBACK_LINK" - iPad fallback link has an
	// invalid (non http/https) URI scheme.
	//   "BAD_DEBUG_PARAM" - Debug param format is incorrect.
	//   "BAD_AD_PARAM" - isAd param format is incorrect.
	//   "DEPRECATED_PARAM" - Indicates a certain param is deprecated.
	//   "UNRECOGNIZED_PARAM" - Indicates certain paramater is not
	// recognized.
	//   "TOO_LONG_PARAM" - Indicates certain paramater is too long.
	//   "NOT_URI_SOCIAL_IMAGE_LINK" - Social meta tag image link is not a
	// valid URI.
	//   "BAD_URI_SCHEME_SOCIAL_IMAGE_LINK" - Social meta tag image link has
	// an invalid (non http/https) URI scheme.
	//   "NOT_URI_SOCIAL_URL"
	//   "BAD_URI_SCHEME_SOCIAL_URL"
	//   "LINK_LENGTH_TOO_LONG" - Dynamic Link URL length is too long.
	//   "LINK_WITH_FRAGMENTS" - Dynamic Link URL contains fragments.
	//   "NOT_MATCHING_IOS_BUNDLE_ID_AND_STORE_ID" - The iOS bundle ID does
	// not match with the given iOS store ID.
	WarningCode string `json:"warningCode,omitempty"`

	// WarningDocumentLink: The document describing the warning, and helps
	// resolve.
	WarningDocumentLink string `json:"warningDocumentLink,omitempty"`

	// WarningMessage: The warning message to help developers improve their
	// requests.
	WarningMessage string `json:"warningMessage,omitempty"`

	// ForceSendFields is a list of field names (e.g. "WarningCode") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "WarningCode") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *DynamicLinkWarning) MarshalJSON() ([]byte, error) {
	type NoMethod DynamicLinkWarning
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// GetIosPostInstallAttributionRequest: Request for iSDK to execute
// strong match flow for post-install attribution.
// This is meant for iOS requests only. Requests from other platforms
// will
// not be honored.
type GetIosPostInstallAttributionRequest struct {
	// AppInstallationTime: App installation epoch time
	// (https://en.wikipedia.org/wiki/Unix_time).
	// This is a client signal for a more accurate weak match.
	AppInstallationTime int64 `json:"appInstallationTime,omitempty,string"`

	// BundleId: APP bundle ID.
	BundleId string `json:"bundleId,omitempty"`

	// Device: Device information.
	Device *DeviceInfo `json:"device,omitempty"`

	// IosVersion: iOS version, ie: 9.3.5.
	// Consider adding "build".
	IosVersion string `json:"iosVersion,omitempty"`

	// RetrievalMethod: App post install attribution retrieval information.
	// Disambiguates
	// mechanism (iSDK or developer invoked) to retrieve payload
	// from
	// clicked link.
	//
	// Possible values:
	//   "UNKNOWN_PAYLOAD_RETRIEVAL_METHOD" - Unknown method.
	//   "IMPLICIT_WEAK_MATCH" - iSDK performs a server lookup by device
	// fingerprint in the background
	// when app is first-opened; no API called by developer.
	//   "EXPLICIT_WEAK_MATCH" - iSDK performs a server lookup by device
	// fingerprint upon a dev API call.
	//   "EXPLICIT_STRONG_AFTER_WEAK_MATCH" - iSDK performs a strong match
	// only if weak match is found upon a dev
	// API call.
	RetrievalMethod string `json:"retrievalMethod,omitempty"`

	// SdkVersion: Google SDK version.
	SdkVersion string `json:"sdkVersion,omitempty"`

	// UniqueMatchLinkToCheck: Possible unique matched link that server need
	// to check before performing
	// fingerprint match. If passed link is short server need to expand the
	// link.
	// If link is long server need to vslidate the link.
	UniqueMatchLinkToCheck string `json:"uniqueMatchLinkToCheck,omitempty"`

	// VisualStyle: Strong match page information. Disambiguates between
	// default UI and
	// custom page to present when strong match succeeds/fails to find
	// cookie.
	//
	// Possible values:
	//   "UNKNOWN_VISUAL_STYLE" - Unknown style.
	//   "DEFAULT_STYLE" - Default style.
	//   "CUSTOM_STYLE" - Custom style.
	VisualStyle string `json:"visualStyle,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AppInstallationTime")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AppInstallationTime") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *GetIosPostInstallAttributionRequest) MarshalJSON() ([]byte, error) {
	type NoMethod GetIosPostInstallAttributionRequest
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// GetIosPostInstallAttributionResponse: Response for iSDK to execute
// strong match flow for post-install attribution.
type GetIosPostInstallAttributionResponse struct {
	// AppMinimumVersion: The minimum version for app, specified by dev
	// through ?imv= parameter.
	// Return to iSDK to allow app to evaluate if current version meets
	// this.
	AppMinimumVersion string `json:"appMinimumVersion,omitempty"`

	// AttributionConfidence: The confidence of the returned attribution.
	//
	// Possible values:
	//   "UNKNOWN_ATTRIBUTION_CONFIDENCE" - Unset.
	//   "WEAK" - Weak confidence, more than one matching link found or link
	// suspected to
	// be false positive
	//   "DEFAULT" - Default confidence, match based on fingerprint
	//   "UNIQUE" - Unique confidence, match based on "unique match link to
	// check" or other
	// means
	AttributionConfidence string `json:"attributionConfidence,omitempty"`

	// DeepLink: The deep-link attributed post-install via one of several
	// techniques
	// (fingerprint, copy unique).
	DeepLink string `json:"deepLink,omitempty"`

	// ExternalBrowserDestinationLink: User-agent specific custom-scheme
	// URIs for iSDK to open. This will be set
	// according to the user-agent tha the click was originally made in.
	// There is
	// no Safari-equivalent custom-scheme open URLs.
	// ie: googlechrome://www.example.com
	// ie: firefox://open-url?url=http://www.example.com
	// ie: opera-http://example.com
	ExternalBrowserDestinationLink string `json:"externalBrowserDestinationLink,omitempty"`

	// FallbackLink: The link to navigate to update the app if min version
	// is not met.
	// This is either (in order): 1) fallback link (from ?ifl= parameter,
	// if
	// specified by developer) or 2) AppStore URL (from ?isi= parameter,
	// if
	// specified), or 3) the payload link (from required link= parameter).
	FallbackLink string `json:"fallbackLink,omitempty"`

	// InvitationId: Invitation ID attributed post-install via one of
	// several techniques
	// (fingerprint, copy unique).
	InvitationId string `json:"invitationId,omitempty"`

	// IsStrongMatchExecutable: Instruction for iSDK to attemmpt to perform
	// strong match. For instance,
	// if browser does not support/allow cookie or outside of support
	// browsers,
	// this will be false.
	IsStrongMatchExecutable bool `json:"isStrongMatchExecutable,omitempty"`

	// MatchMessage: Describes why match failed, ie: "discarded due to low
	// confidence".
	// This message will be publicly visible.
	MatchMessage string `json:"matchMessage,omitempty"`

	// RequestedLink: Entire FDL (short or long) attributed post-install via
	// one of several
	// techniques (fingerprint, copy unique).
	RequestedLink string `json:"requestedLink,omitempty"`

	// ResolvedLink: The entire FDL, expanded from a short link. It is the
	// same as the
	// requested_link, if it is long. Parameters from this should not
	// be
	// used directly (ie: server can default utm_[campaign|medium|source]
	// to a value when requested_link lack them, server determine the
	// best
	// fallback_link when requested_link specifies >1 fallback links).
	ResolvedLink string `json:"resolvedLink,omitempty"`

	// UtmCampaign: Scion campaign value to be propagated by iSDK to Scion
	// at post-install.
	UtmCampaign string `json:"utmCampaign,omitempty"`

	// UtmMedium: Scion medium value to be propagated by iSDK to Scion at
	// post-install.
	UtmMedium string `json:"utmMedium,omitempty"`

	// UtmSource: Scion source value to be propagated by iSDK to Scion at
	// post-install.
	UtmSource string `json:"utmSource,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "AppMinimumVersion")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AppMinimumVersion") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *GetIosPostInstallAttributionResponse) MarshalJSON() ([]byte, error) {
	type NoMethod GetIosPostInstallAttributionResponse
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// GooglePlayAnalytics: Parameters for Google Play Campaign
// Measurements.
// [Learn
// more](https://developers.google.com/analytics/devguides/collection/and
// roid/v4/campaigns#campaign-params)
type GooglePlayAnalytics struct {
	// Gclid: [AdWords autotagging
	// parameter](https://support.google.com/analytics/answer/1033981?hl=en);
	//
	// used to measure Google AdWords ads. This value is generated
	// dynamically
	// and should never be modified.
	Gclid string `json:"gclid,omitempty"`

	// UtmCampaign: Campaign name; used for keyword analysis to identify a
	// specific product
	// promotion or strategic campaign.
	UtmCampaign string `json:"utmCampaign,omitempty"`

	// UtmContent: Campaign content; used for A/B testing and
	// content-targeted ads to
	// differentiate ads or links that point to the same URL.
	UtmContent string `json:"utmContent,omitempty"`

	// UtmMedium: Campaign medium; used to identify a medium such as email
	// or cost-per-click.
	UtmMedium string `json:"utmMedium,omitempty"`

	// UtmSource: Campaign source; used to identify a search engine,
	// newsletter, or other
	// source.
	UtmSource string `json:"utmSource,omitempty"`

	// UtmTerm: Campaign term; used with paid search to supply the keywords
	// for ads.
	UtmTerm string `json:"utmTerm,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Gclid") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Gclid") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *GooglePlayAnalytics) MarshalJSON() ([]byte, error) {
	type NoMethod GooglePlayAnalytics
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ITunesConnectAnalytics: Parameters for iTunes Connect App Analytics.
type ITunesConnectAnalytics struct {
	// At: Affiliate token used to create affiliate-coded links.
	At string `json:"at,omitempty"`

	// Ct: Campaign text that developers can optionally add to any link in
	// order to
	// track sales from a specific marketing campaign.
	Ct string `json:"ct,omitempty"`

	// Mt: iTune media types, including music, podcasts, audiobooks and so
	// on.
	Mt string `json:"mt,omitempty"`

	// Pt: Provider token that enables analytics for Dynamic Links from
	// within iTunes
	// Connect.
	Pt string `json:"pt,omitempty"`

	// ForceSendFields is a list of field names (e.g. "At") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "At") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ITunesConnectAnalytics) MarshalJSON() ([]byte, error) {
	type NoMethod ITunesConnectAnalytics
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// IosInfo: iOS related attributes to the Dynamic Link..
type IosInfo struct {
	// IosAppStoreId: iOS App Store ID.
	IosAppStoreId string `json:"iosAppStoreId,omitempty"`

	// IosBundleId: iOS bundle ID of the app.
	IosBundleId string `json:"iosBundleId,omitempty"`

	// IosCustomScheme: Custom (destination) scheme to use for iOS. By
	// default, we’ll use the
	// bundle ID as the custom scheme. Developer can override this behavior
	// using
	// this param.
	IosCustomScheme string `json:"iosCustomScheme,omitempty"`

	// IosFallbackLink: Link to open on iOS if the app is not installed.
	IosFallbackLink string `json:"iosFallbackLink,omitempty"`

	// IosIpadBundleId: iPad bundle ID of the app.
	IosIpadBundleId string `json:"iosIpadBundleId,omitempty"`

	// IosIpadFallbackLink: If specified, this overrides the
	// ios_fallback_link value on iPads.
	IosIpadFallbackLink string `json:"iosIpadFallbackLink,omitempty"`

	// ForceSendFields is a list of field names (e.g. "IosAppStoreId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "IosAppStoreId") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *IosInfo) MarshalJSON() ([]byte, error) {
	type NoMethod IosInfo
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ManagedShortLink: Managed Short Link.
type ManagedShortLink struct {
	// CreationTime: Creation timestamp of the short link.
	CreationTime string `json:"creationTime,omitempty"`

	// FlaggedAttribute: Attributes that have been flagged about this short
	// url.
	//
	// Possible values:
	//   "UNSPECIFIED_ATTRIBUTE" - Indicates that no attributes were found
	// for this short url.
	//   "SPAM" - Indicates that short url has been flagged by AbuseIAm team
	// as spam.
	FlaggedAttribute []string `json:"flaggedAttribute,omitempty"`

	// Info: Full Dyamic Link info
	Info *DynamicLinkInfo `json:"info,omitempty"`

	// Link: Short durable link url, for example,
	// "https://sample.app.goo.gl/xyz123".
	//
	// Required.
	Link string `json:"link,omitempty"`

	// LinkName: Link name defined by the creator.
	//
	// Required.
	LinkName string `json:"linkName,omitempty"`

	// Visibility: Visibility status of link.
	//
	// Possible values:
	//   "UNSPECIFIED_VISIBILITY" - Visibility of the link is not specified.
	//   "UNARCHIVED" - Link created in console and should be shown in
	// console.
	//   "ARCHIVED" - Link created in console and should not be shown in
	// console (but can
	// be shown in the console again if it is unarchived).
	//   "NEVER_SHOWN" - Link created outside of console and should never be
	// shown in console.
	Visibility string `json:"visibility,omitempty"`

	// ForceSendFields is a list of field names (e.g. "CreationTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CreationTime") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ManagedShortLink) MarshalJSON() ([]byte, error) {
	type NoMethod ManagedShortLink
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// NavigationInfo: Information of navigation behavior.
type NavigationInfo struct {
	// EnableForcedRedirect: If this option is on, FDL click will be forced
	// to redirect rather than
	// show an interstitial page.
	EnableForcedRedirect bool `json:"enableForcedRedirect,omitempty"`

	// ForceSendFields is a list of field names (e.g.
	// "EnableForcedRedirect") to unconditionally include in API requests.
	// By default, fields with empty values are omitted from API requests.
	// However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "EnableForcedRedirect") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *NavigationInfo) MarshalJSON() ([]byte, error) {
	type NoMethod NavigationInfo
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// SocialMetaTagInfo: Parameters for social meta tag params.
// Used to set meta tag data for link previews on social sites.
type SocialMetaTagInfo struct {
	// SocialDescription: A short description of the link. Optional.
	SocialDescription string `json:"socialDescription,omitempty"`

	// SocialImageLink: An image url string. Optional.
	SocialImageLink string `json:"socialImageLink,omitempty"`

	// SocialTitle: Title to be displayed. Optional.
	SocialTitle string `json:"socialTitle,omitempty"`

	// ForceSendFields is a list of field names (e.g. "SocialDescription")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "SocialDescription") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *SocialMetaTagInfo) MarshalJSON() ([]byte, error) {
	type NoMethod SocialMetaTagInfo
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Suffix: Short Dynamic Link suffix.
type Suffix struct {
	// CustomSuffix: Only applies to Option.CUSTOM.
	CustomSuffix string `json:"customSuffix,omitempty"`

	// Option: Suffix option.
	//
	// Possible values:
	//   "OPTION_UNSPECIFIED" - The suffix option is not specified, performs
	// as UNGUESSABLE .
	//   "UNGUESSABLE" - Short Dynamic Link suffix is a base62 [0-9A-Za-z]
	// encoded string of
	// a random generated 96 bit random number, which has a length of 17
	// chars.
	// For example, "nlAR8U4SlKRZw1cb2".
	// It prevents other people from guessing and crawling short Dynamic
	// Links
	// that contain personal identifiable information.
	//   "SHORT" - Short Dynamic Link suffix is a base62 [0-9A-Za-z] string
	// starting with a
	// length of 4 chars. the length will increase when all the space
	// is
	// occupied.
	//   "CUSTOM" - Custom DDL suffix is a client specified string, for
	// example,
	// "buy2get1free".
	// NOTE: custom suffix should only be available to managed short
	// link
	// creation
	Option string `json:"option,omitempty"`

	// ForceSendFields is a list of field names (e.g. "CustomSuffix") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CustomSuffix") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Suffix) MarshalJSON() ([]byte, error) {
	type NoMethod Suffix
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// method id "firebasedynamiclinks.managedShortLinks.create":

type ManagedShortLinksCreateCall struct {
	s                             *Service
	createmanagedshortlinkrequest *CreateManagedShortLinkRequest
	urlParams_                    gensupport.URLParams
	ctx_                          context.Context
	header_                       http.Header
}

// Create: Creates a managed short Dynamic Link given either a valid
// long Dynamic Link
// or details such as Dynamic Link domain, Android and iOS app
// information.
// The created short Dynamic Link will not expire.
//
// This differs from CreateShortDynamicLink in the following ways:
//   - The request will also contain a name for the link (non unique
// name
//     for the front end).
//   - The response must be authenticated with an auth token (generated
// with
//     the admin service account).
//   - The link will appear in the FDL list of links in the console
// front end.
//
// The Dynamic Link domain in the request must be owned by
// requester's
// Firebase project.
func (r *ManagedShortLinksService) Create(createmanagedshortlinkrequest *CreateManagedShortLinkRequest) *ManagedShortLinksCreateCall {
	c := &ManagedShortLinksCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.createmanagedshortlinkrequest = createmanagedshortlinkrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ManagedShortLinksCreateCall) Fields(s ...googleapi.Field) *ManagedShortLinksCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ManagedShortLinksCreateCall) Context(ctx context.Context) *ManagedShortLinksCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ManagedShortLinksCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ManagedShortLinksCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.createmanagedshortlinkrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/managedShortLinks:create")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "firebasedynamiclinks.managedShortLinks.create" call.
// Exactly one of *CreateManagedShortLinkResponse or error will be
// non-nil. Any non-2xx status code is an error. Response headers are in
// either *CreateManagedShortLinkResponse.ServerResponse.Header or (if a
// response was returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ManagedShortLinksCreateCall) Do(opts ...googleapi.CallOption) (*CreateManagedShortLinkResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &CreateManagedShortLinkResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Creates a managed short Dynamic Link given either a valid long Dynamic Link\nor details such as Dynamic Link domain, Android and iOS app information.\nThe created short Dynamic Link will not expire.\n\nThis differs from CreateShortDynamicLink in the following ways:\n  - The request will also contain a name for the link (non unique name\n    for the front end).\n  - The response must be authenticated with an auth token (generated with\n    the admin service account).\n  - The link will appear in the FDL list of links in the console front end.\n\nThe Dynamic Link domain in the request must be owned by requester's\nFirebase project.",
	//   "flatPath": "v1/managedShortLinks:create",
	//   "httpMethod": "POST",
	//   "id": "firebasedynamiclinks.managedShortLinks.create",
	//   "parameterOrder": [],
	//   "parameters": {},
	//   "path": "v1/managedShortLinks:create",
	//   "request": {
	//     "$ref": "CreateManagedShortLinkRequest"
	//   },
	//   "response": {
	//     "$ref": "CreateManagedShortLinkResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/firebase"
	//   ]
	// }

}

// method id "firebasedynamiclinks.shortLinks.create":

type ShortLinksCreateCall struct {
	s                             *Service
	createshortdynamiclinkrequest *CreateShortDynamicLinkRequest
	urlParams_                    gensupport.URLParams
	ctx_                          context.Context
	header_                       http.Header
}

// Create: Creates a short Dynamic Link given either a valid long
// Dynamic Link or
// details such as Dynamic Link domain, Android and iOS app
// information.
// The created short Dynamic Link will not expire.
//
// Repeated calls with the same long Dynamic Link or Dynamic Link
// information
// will produce the same short Dynamic Link.
//
// The Dynamic Link domain in the request must be owned by
// requester's
// Firebase project.
func (r *ShortLinksService) Create(createshortdynamiclinkrequest *CreateShortDynamicLinkRequest) *ShortLinksCreateCall {
	c := &ShortLinksCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.createshortdynamiclinkrequest = createshortdynamiclinkrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ShortLinksCreateCall) Fields(s ...googleapi.Field) *ShortLinksCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ShortLinksCreateCall) Context(ctx context.Context) *ShortLinksCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ShortLinksCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ShortLinksCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.createshortdynamiclinkrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/shortLinks")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "firebasedynamiclinks.shortLinks.create" call.
// Exactly one of *CreateShortDynamicLinkResponse or error will be
// non-nil. Any non-2xx status code is an error. Response headers are in
// either *CreateShortDynamicLinkResponse.ServerResponse.Header or (if a
// response was returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ShortLinksCreateCall) Do(opts ...googleapi.CallOption) (*CreateShortDynamicLinkResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &CreateShortDynamicLinkResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Creates a short Dynamic Link given either a valid long Dynamic Link or\ndetails such as Dynamic Link domain, Android and iOS app information.\nThe created short Dynamic Link will not expire.\n\nRepeated calls with the same long Dynamic Link or Dynamic Link information\nwill produce the same short Dynamic Link.\n\nThe Dynamic Link domain in the request must be owned by requester's\nFirebase project.",
	//   "flatPath": "v1/shortLinks",
	//   "httpMethod": "POST",
	//   "id": "firebasedynamiclinks.shortLinks.create",
	//   "parameterOrder": [],
	//   "parameters": {},
	//   "path": "v1/shortLinks",
	//   "request": {
	//     "$ref": "CreateShortDynamicLinkRequest"
	//   },
	//   "response": {
	//     "$ref": "CreateShortDynamicLinkResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/firebase"
	//   ]
	// }

}

// method id "firebasedynamiclinks.getLinkStats":

type V1GetLinkStatsCall struct {
	s            *Service
	dynamicLink  string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// GetLinkStats: Fetches analytics stats of a short Dynamic Link for a
// given
// duration. Metrics include number of clicks, redirects, installs,
// app first opens, and app reopens.
func (r *V1Service) GetLinkStats(dynamicLink string) *V1GetLinkStatsCall {
	c := &V1GetLinkStatsCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.dynamicLink = dynamicLink
	return c
}

// DurationDays sets the optional parameter "durationDays": The span of
// time requested in days.
func (c *V1GetLinkStatsCall) DurationDays(durationDays int64) *V1GetLinkStatsCall {
	c.urlParams_.Set("durationDays", fmt.Sprint(durationDays))
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *V1GetLinkStatsCall) Fields(s ...googleapi.Field) *V1GetLinkStatsCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *V1GetLinkStatsCall) IfNoneMatch(entityTag string) *V1GetLinkStatsCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *V1GetLinkStatsCall) Context(ctx context.Context) *V1GetLinkStatsCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *V1GetLinkStatsCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *V1GetLinkStatsCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/{dynamicLink}/linkStats")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"dynamicLink": c.dynamicLink,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "firebasedynamiclinks.getLinkStats" call.
// Exactly one of *DynamicLinkStats or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *DynamicLinkStats.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *V1GetLinkStatsCall) Do(opts ...googleapi.CallOption) (*DynamicLinkStats, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &DynamicLinkStats{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Fetches analytics stats of a short Dynamic Link for a given\nduration. Metrics include number of clicks, redirects, installs,\napp first opens, and app reopens.",
	//   "flatPath": "v1/{dynamicLink}/linkStats",
	//   "httpMethod": "GET",
	//   "id": "firebasedynamiclinks.getLinkStats",
	//   "parameterOrder": [
	//     "dynamicLink"
	//   ],
	//   "parameters": {
	//     "durationDays": {
	//       "description": "The span of time requested in days.",
	//       "format": "int64",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "dynamicLink": {
	//       "description": "Dynamic Link URL. e.g. https://abcd.app.goo.gl/wxyz",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/{dynamicLink}/linkStats",
	//   "response": {
	//     "$ref": "DynamicLinkStats"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/firebase"
	//   ]
	// }

}

// method id "firebasedynamiclinks.installAttribution":

type V1InstallAttributionCall struct {
	s                                   *Service
	getiospostinstallattributionrequest *GetIosPostInstallAttributionRequest
	urlParams_                          gensupport.URLParams
	ctx_                                context.Context
	header_                             http.Header
}

// InstallAttribution: Get iOS strong/weak-match info for post-install
// attribution.
func (r *V1Service) InstallAttribution(getiospostinstallattributionrequest *GetIosPostInstallAttributionRequest) *V1InstallAttributionCall {
	c := &V1InstallAttributionCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.getiospostinstallattributionrequest = getiospostinstallattributionrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *V1InstallAttributionCall) Fields(s ...googleapi.Field) *V1InstallAttributionCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *V1InstallAttributionCall) Context(ctx context.Context) *V1InstallAttributionCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *V1InstallAttributionCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *V1InstallAttributionCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.getiospostinstallattributionrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/installAttribution")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "firebasedynamiclinks.installAttribution" call.
// Exactly one of *GetIosPostInstallAttributionResponse or error will be
// non-nil. Any non-2xx status code is an error. Response headers are in
// either *GetIosPostInstallAttributionResponse.ServerResponse.Header or
// (if a response was returned at all) in
// error.(*googleapi.Error).Header. Use googleapi.IsNotModified to check
// whether the returned error was because http.StatusNotModified was
// returned.
func (c *V1InstallAttributionCall) Do(opts ...googleapi.CallOption) (*GetIosPostInstallAttributionResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &GetIosPostInstallAttributionResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Get iOS strong/weak-match info for post-install attribution.",
	//   "flatPath": "v1/installAttribution",
	//   "httpMethod": "POST",
	//   "id": "firebasedynamiclinks.installAttribution",
	//   "parameterOrder": [],
	//   "parameters": {},
	//   "path": "v1/installAttribution",
	//   "request": {
	//     "$ref": "GetIosPostInstallAttributionRequest"
	//   },
	//   "response": {
	//     "$ref": "GetIosPostInstallAttributionResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/firebase"
	//   ]
	// }

}
