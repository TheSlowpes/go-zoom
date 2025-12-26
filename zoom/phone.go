package zoom

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

type PhoneService struct {
	client            *Client
	Accounts          *PhoneAccountsService
	Alerts            *PhoneAlertsService
	AudioLibrary      *PhoneAudioLibraryService
	AutoReceptionists *PhoneAutoReceptionistsService
	BillingAccounts   *PhoneBillingAccountService
	BlockedLists      *PhoneBlockedListService
	CallHandling      *PhoneCallHandlingService
}

type PhoneAccountsService struct {
	client *Client
}

type AddCustomizedNumbersRequest struct {
	PhoneNumberIDs []string `json:"phone_number_ids"`
}

// https://developers.zoom.us/docs/api/phone/#tag/accounts/post/phone/outbound_caller_id/customized_numbers
func (p *PhoneAccountsService) AddCustomizedNumbers(ctx context.Context, req *AddCustomizedNumbersRequest) (*http.Response, error) {
	if len(req.PhoneNumberIDs) > 30 {
		return nil, fmt.Errorf("Error: cannot add more than 30 phone numbers ids at once")
	}
	res, err := p.client.request(ctx, http.MethodPost, "/phone/outbound_caller_id/customized_numbers", nil, req, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type DeleteCustomizedNumbersRequest struct {
	CustomizedIDs []string `url:"customized_ids"`
}

// https://developers.zoom.us/docs/api/phone/#tag/accounts/delete/phone/outbound_caller_id/customized_numbers
func (p *PhoneAccountsService) DeleteCustomizedNumbers(ctx context.Context, req *DeleteCustomizedNumbersRequest) (*http.Response, error) {
	if len(req.CustomizedIDs) > 30 {
		return nil, fmt.Errorf("Error: cannot delete more than 30 customized ids at once")
	}
	res, err := p.client.request(ctx, http.MethodDelete, "/phone/outbound_caller_id/customized_numbers", req, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type GetCustomizedNumbersRequest struct {
	*PaginationOptions `url:",omitempty"`
}

type CustomizeNumber struct {
	CustomizeID     string `json:"customize_id"`
	DisplayName     string `json:"display_name"`
	ExtensionID     string `json:"extension_id"`
	ExtensionName   string `json:"extension_name"`
	ExtensionNumber string `json:"extension_number"`
	ExtensionType   string `json:"extension_type"`
	Incoming        bool   `json:"incoming"`
	Outgoing        bool   `json:"outgoing"`
	PhoneNumber     string `json:"phone_number"`
	PhoneNumberID   string `json:"phone_number_id"`
	Site            struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"site"`
}

type GetCustomizedNumbersResponse struct {
	*PaginationResponse
	CustomizeNumbers []*CustomizeNumber `json:"customize_numbers"`
}

// https://developers.zoom.us/docs/api/phone/#tag/accounts/get/phone/outbound_caller_id/customized_numbers
func (p *PhoneAccountsService) GetCustomizedNumbers(ctx context.Context, req *GetCustomizedNumbersRequest) (*GetCustomizedNumbersResponse, *http.Response, error) {
	out := &GetCustomizedNumbersResponse{}

	res, err := p.client.request(ctx, http.MethodGet, "/phone/outbound_caller_id/customized_numbers", req, nil, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

var availableSettingTypes = []string{
	"call_live_transcription",
	"local_survivability_mode",
	"external_calling_on_zoom_room_common_area",
	"select_outbound_caller_id",
	"personal_audio_library",
	"voicemail",
	"voicemail_transcription",
	"voicemail_notification_by_email",
	"shared_voicemail_notification_by_email",
	"restricted_call_hours",
	"allowed_call_locations",
	"check_voicemails_over_phone",
	"auto_call_recording",
	"ad_hoc_call_recording",
	"international_calling",
	"outbound_calling",
	"outbound_sms",
	"sms",
	"sms_etiquette_tool",
	"zoom_phone_on_mobile",
	"zoom_phone_on_pwa",
	"e2e_encryption",
	"call_handling_forwarding_to_other_users",
	"call_overflow",
	"call_transferring",
	"elevate_to_meeting",
	"call_park",
	"hand_off_to_room",
	"mobile_switch_to_carrier",
	"delegation",
	"audio_intercom",
	"block_calls_without_caller_id",
	"block_external_calls",
	"call_queue_opt_out_reason",
	"auto_delete_data_after_retention_duration",
	"auto_call_from_third_party_apps",
	"override_default_port",
	"peer_to_peer_media",
	"advanced_encryption",
	"display_call_feedback_survey",
	"block_list_for_inbound_calls_and_messaging",
	"block_calls_as_threat",
}

type AccountSettingsQuery struct {
	SettingTypes string `url:"setting_type"`
}

type AccountSettingStates struct {
	Enable   bool   `json:"enable"`
	Locked   bool   `json:"locked"`
	LockedBy string `json:"locked_by"` //invalid, account
}

type AccountSettingsResponse struct {
	AdHocCallRecording struct {
		*AccountSettingStates
	} `json:"ad_hoc_call_recording"`
	AdavancedEncryption struct {
		*AccountSettingStates
	} `json:"advanced_encryption"`
	AllowedCallLocations struct {
		AllowInternalCalls bool `json:"allow_internal_calls"`
		*AccountSettingStates
	} `json:"allowed_call_locations"`
	AudioIntercom struct {
		*AccountSettingStates
	} `json:"audio_intercom"`
	AutoCallFromThirdPartyApps struct {
		*AccountSettingStates
	} `json:"auto_call_from_third_party_apps"`
	AutoCallRecording struct {
		*AccountSettingStates
		AllowStopResumeRecording     bool `json:"allow_stop_resume_recording"`
		DisconnectOnRecordingFailure bool `json:"disconnect_on_recording_failure"`
		InboundAudioNotification     struct {
			RecordingExplicitConsent    bool   `json:"recording_explicit_consent"`
			RecordingStartPrompt        bool   `json:"recording_start_prompt"`
			RecordingStartPromptAudioID string `json:"recording_start_prompt_audio_id"`
		} `json:"inbound_audio_notification"`
		OutboundAudioNotification struct {
			RecordingExplicitConsent    bool   `json:"recording_explicit_consent"`
			RecordingStartPrompt        bool   `json:"recording_start_prompt"`
			RecordingStartPromptAudioID string `json:"recording_start_prompt_audio_id"`
		} `json:"outbound_audio_notification"`
		PlayRecordingBeepTone struct {
			Enable               bool   `json:"enable"`
			PlayBeepMember       string `json:"play_beep_member"`
			PlayBeepTimeInterval int    `json:"play_beep_time_interval"`
		} `json:"play_recording_beep_tone"`
		RecordingCalls         string `json:"recording_calls"`
		RecordingTranscription bool   `json:"recording_transcription"`
	} `json:"auto_call_recording"`
	AutoDeleteDataAfterRetentionDuration struct {
		*AccountSettingStates
	} `json:"auto_delete_data_after_retention_duration"`
	BlockCallsAsThreat struct {
		*AccountSettingStates
	} `json:"block_calls_as_threat"`
	BlockCallsWithoutCallerID struct {
		*AccountSettingStates
	} `json:"block_calls_without_caller_id"`
	BlockExternalCalls struct {
		*AccountSettingStates
	} `json:"block_external_calls"`
	BlockListForInboundCallsAndMessaging struct {
		*AccountSettingStates
	} `json:"block_list_for_inbound_calls_and_messaging"`
	CallHandlingForwardingToOtherUsers struct {
		*AccountSettingStates
		CallForwardingType int `json:"call_forwarding_type"`
	} `json:"call_handling_forwarding_to_other_users"`
	CallLiveTranscription struct {
		*AccountSettingStates
		TranscriptionStartPrompt struct {
			AudioID   string `json:"audio_id"`
			AudioName string `json:"audio_name"`
		} `json:"transcription_start_prompt"`
	} `json:"call_live_transcription"`
	CallOverflow struct {
		*AccountSettingStates
		CallOverflowType int `json:"call_overflow_type"`
	} `json:"call_overflow"`
	CallPark struct {
		*AccountSettingStates
	} `json:"call_park"`
	CallQueueOptOutReason struct {
		*AccountSettingStates
	} `json:"call_queue_opt_out_reason"`
	CallTransferring struct {
		*AccountSettingStates
		CallTransferringType int `json:"call_transferring_type"`
	} `json:"call_transferring"`
	CheckVoicemailsOverPhone struct {
		*AccountSettingStates
	} `json:"check_voicemails_over_phone"`
	Delegation struct {
		*AccountSettingStates
	} `json:"delegation"`
	DisplayCallFeedbackSurvey struct {
		*AccountSettingStates
	} `json:"display_call_feedback_survey"`
	E2EEncryption struct {
		*AccountSettingStates
	} `json:"e2e_encryption"`
	ElevateToMeeting struct {
		*AccountSettingStates
	} `json:"elevate_to_meeting"`
	ExternalCallingOnZoomRoomCommonArea struct {
		*AccountSettingStates
	} `json:"external_calling_on_zoom_room_common_area"`
	HandOffToRoom struct {
		*AccountSettingStates
	} `json:"hand_off_to_room"`
	InternationalCalling struct {
		*AccountSettingStates
	} `json:"international_calling"`
	LocalSurvivabilityMode struct {
		*AccountSettingStates
	} `json:"local_survivability_mode"`
	MobileSwitchToCarrier struct {
		*AccountSettingStates
	} `json:"mobile_switch_to_carrier"`
	OutboundCalling struct {
		*AccountSettingStates
	} `json:"outbound_calling"`
	OutboundSMS struct {
		*AccountSettingStates
	} `json:"outbound_sms"`
	OverrideDefaultPort struct {
		*AccountSettingStates
	} `json:"override_default_port"`
	PeerToPeerMedia struct {
		*AccountSettingStates
	} `json:"peer_to_peer_media"`
	PersonalAudioLibrary struct {
		*AccountSettingStates
		AllowMusicOnHoldCustomization                 bool `json:"allow_music_on_hold_customization"`
		AllowVoicemailAndMessageGreetingCustomization bool `json:"allow_voicemail_and_message_greeting_customization"`
	} `json:"restricted_call_hours"`
	RestrictedCallHours struct {
		*AccountSettingStates
	} `json:"personal_audio_library"`
}

// https://developers.zoom.us/docs/api/phone/#tag/accounts/get/phone/account_settings
func (p *PhoneAccountsService) GetAccountSettings(ctx context.Context, query *AccountSettingsQuery) (*AccountSettingsResponse, *http.Response, error) {
	for _, setting := range strings.Split(query.SettingTypes, ",") {
		setting = strings.TrimSpace(setting)
		if !slices.Contains(availableSettingTypes, setting) {
			return nil, nil, fmt.Errorf("Error: invalid setting type '%s'", setting)
		}
	}
	out := &AccountSettingsResponse{}

	res, err := p.client.request(ctx, http.MethodGet, "/phone/account_settings", query, nil, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type PhoneAlertsService struct {
	client *Client
}

type CreateAlertRequest struct {
	AlertSettingsName string `json:"alert_settings_name"`
	Module            int    `json:"module"`
	Rule              int    `json:"rule"`
	RuleConditions    []struct {
		RuleConditionType  int    `json:"rule_condition_type"`
		RuleConditionValue string `json:"rule_condition_value"`
	} `json:"rule_conditions"`
	TargetType    int    `json:"target_type"`
	TimeFrameFrom string `json:"time_frame_from"`
	TimeFrameTo   string `json:"time_frame_to"`
	TimeFrameType string `json:"time_frame_type"`
	ChatChannels  []struct {
		ChatChannelName string `json:"chat_channel_name"`
		EndPoint        string `json:"endpoint"`
		Token           string `json:"token"`
	} `json:"chat_channels,omitempty"`
	EmailRecipients []string `json:"email_recipients,omitempty"`
	Frequency       int      `json:"frequency"`
	TargetIDs       []string `json:"target_ids"`
	Status          int      `json:"status"`
}

type CreateAlertResponse struct {
	AlertSettingID   string `json:"alert_setting_id"`
	AlertSettingName string `json:"alert_setting_name"`
}

// https://developers.zoom.us/docs/api/phone/#tag/alerts/post/phone/alert_settings
func (p *PhoneAlertsService) CreateAlert(ctx context.Context, req *CreateAlertRequest) (*CreateAlertResponse, *http.Response, error) {
	out := &CreateAlertResponse{}

	res, err := p.client.request(ctx, http.MethodPost, "/phone/alert_settings", nil, req, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type DeleteAlertRequest struct {
	AlertSettingID string
}

// https://developers.zoom.us/docs/api/phone/#tag/alerts/delete/phone/alert_settings/%7BalertSettingId%7D
func (p *PhoneAlertsService) DeleteAlert(ctx context.Context, req *DeleteAlertRequest) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodDelete, fmt.Sprintf("/phone/alert_settings/%s", req.AlertSettingID), nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type GetAlertSettingsRequest struct {
	AlertSettingID string
}

type GetAlertSettingsResponse struct {
	AlertSettingID   string `json:"alert_setting_id"`
	AlertSettingName string `json:"alert_setting_name"`
	ChatChannels     []struct {
		ChatChannelName string `json:"chat_channel_name"`
		EndPoint        string `json:"endpoint"`
		Token           string `json:"token"`
	} `json:"chat_channels"`
	EmailRecipients []string `json:"email_recipients"`
	Frequency       int      `json:"frequency"`
	Module          int      `json:"module"`
	Rule            int      `json:"rule"`
	RuleConditions  []struct {
		RuleConditionType  int    `json:"rule_condition_type"`
		RuleConditionValue string `json:"rule_condition_value"`
	} `json:"rule_conditions"`
}

// https://developers.zoom.us/docs/api/phone/#tag/alerts/get/phone/alert_settings/{alertSettingId}
func (p *PhoneAlertsService) GetAlertSettingsDetails(ctx context.Context, req *GetAlertSettingsRequest) (*GetAlertSettingsResponse, *http.Response, error) {
	out := &GetAlertSettingsResponse{}
	res, err := p.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/alert_settings/%s", req.AlertSettingID), nil, nil, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type ListAlertSettingsRequest struct {
	*PaginationOptions
	Module int `url:"module"`
	Rule   int `url:"rule"`
	Status int `url:"status"`
}

type ListAlertSettingsResponse struct {
	*PaginationResponse
	AlertSettings []*GetAlertSettingsResponse
}

// https://developers.zoom.us/docs/api/phone/#tag/alerts/get/phone/alert_settings
func (p *PhoneAlertsService) ListAlertSettings(ctx context.Context, req *ListAlertSettingsRequest) (*ListAlertSettingsResponse, *http.Response, error) {
	out := &ListAlertSettingsResponse{}
	res, err := p.client.request(ctx, http.MethodGet, "/phone/alert_settings", req, nil, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type UpdateAlertSettingsRequest struct {
	AlertSettingsID string
}

// https://developers.zoom.us/docs/api/phone/#tag/alerts/patch/phone/alert_settings/{alertSettingId}
func (p *PhoneAlertsService) UpdateAlertSettings(ctx context.Context, req *UpdateAlertSettingsRequest, body *CreateAlertRequest) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodPatch, fmt.Sprintf("/phone/alert_settings/%s", req.AlertSettingsID), nil, body, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type PhoneAudioLibraryService struct {
	client *Client
}

type AddAudioItemPathParam struct {
	UserId string
}

type AddAudioItemRequest struct {
	AudioName     string `json:"audio_name"`
	Text          string `json:"text,omitempty"`
	VoiceAccent   string `json:"voice_accent,omitempty"`
	VoiceLanguage string `json:"voice_language,omitempty"`
}

type AddAudioItemResponse struct {
	AudioID string `json:"audio_id"`
	Name    string `json:"name"`
}

// https://developers.zoom.us/docs/api/phone/#tag/audio-library/post/phone/users/{userId}/audios
func (p *PhoneAudioLibraryService) AddAudioItem(ctx context.Context, pathParam *AddAudioItemPathParam, req *AddAudioItemRequest) (*AddAudioItemResponse, *http.Response, error) {
	out := &AddAudioItemResponse{}

	res, err := p.client.request(ctx, http.MethodPost, fmt.Sprintf("/phone/users/%s/audios", pathParam.UserId), nil, req, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type AddAudioItemsPathParam struct {
	UserId string
}

type AddAudioItemsRequest struct {
	Attachments []struct {
		AudioType       string `json:"audio_type"`
		Base64Enconding string `json:"base64_encoding"`
		Name            string `json:"name"`
	} `json:"attachments"`
}

type AddAudioItemsResponse struct {
	Audios []*AddAudioItemResponse `json:"audios"`
}

// https://developers.zoom.us/docs/api/phone/#tag/audio-library/post/phone/users/{userId}/audios/batch
func (p *PhoneAudioLibraryService) AddAudioItems(ctx context.Context, pathParam *AddAudioItemsPathParam, req *AddAudioItemsRequest) (*AddAudioItemsResponse, *http.Response, error) {
	out := &AddAudioItemsResponse{}

	res, err := p.client.request(ctx, http.MethodPost, fmt.Sprintf("/phone/users/%s/audios/batch", pathParam.UserId), nil, req, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type DeleteAudioItemPathParam struct {
	AudioID string
}

// https://developers.zoom.us/docs/api/phone/#tag/audio-library/delete/phone/audios/{audioId}
func (p *PhoneAudioLibraryService) DeleteAudioItem(ctx context.Context, pathParam *DeleteAudioItemPathParam) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodDelete, fmt.Sprintf("/phone/audios/%s", pathParam.AudioID), nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type GetAudioItemPathParam struct {
	AudioID string
}

type GetAudioItemResponse struct {
	AudioID       string `json:"audio_id"`
	Name          string `json:"name"`
	PlayURL       string `json:"play_url"`
	Text          string `json:"text,omitempty"`
	VoiceAccent   string `json:"voice_accent,omitempty"`
	VoiceLanguage string `json:"voice_language,omitempty"`
}

// https://developers.zoom.us/docs/api/phone/#tag/audio-library/get/phone/audios/{audioId}
func (p *PhoneAudioLibraryService) GetAudioItem(ctx context.Context, pathParam *GetAudioItemPathParam) (*GetAudioItemResponse, *http.Response, error) {
	out := &GetAudioItemResponse{}

	res, err := p.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/audios/%s", pathParam.AudioID), nil, nil, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type ListAudioItemsPathParam struct {
	UserId string
}

type ListAudioItemsResponse struct {
	Audios []*GetAudioItemResponse `json:"audios"`
}

// https://developers.zoom.us/docs/api/phone/#tag/audio-library/get/phone/users/{userId}/audios
func (p *PhoneAudioLibraryService) ListAudioItems(ctx context.Context, pathParam *ListAudioItemsPathParam) (*ListAudioItemsResponse, *http.Response, error) {
	out := &ListAudioItemsResponse{}

	res, err := p.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/users/%s/audios", pathParam.UserId), nil, nil, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type UpdateAudioItemPathParam struct {
	AudioID string
}

type UpdateAudioItemRequest struct {
	Name string `json:"name"`
}

// https://developers.zoom.us/docs/api/phone/#tag/audio-library/patch/phone/audios/{audioId}
func (p *PhoneAudioLibraryService) UpdateAudioItem(ctx context.Context, pathParam *UpdateAudioItemPathParam, req *UpdateAudioItemRequest) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodPatch, fmt.Sprintf("/phone/audios/%s", pathParam.AudioID), nil, req, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type PhoneAutoReceptionistsService struct {
	client *Client
}

type AddPolicySubSettingPathParams struct {
	autoReceptionistId string
	policyType         string
}

type AutoReceptionistPolicySubSetting struct {
	VoiceMailAccessMember struct {
		AccessUserID   string `json:"access_user_id"`
		AccessUserType string `json:"access_user_type"`
		Delete         bool   `json:"delete"`
		Download       bool   `json:"download"`
	} `json:"voice_mail_access_member"`
}

// https://developers.zoom.us/docs/api/phone/#tag/auto-receptionists/post/phone/auto_receptionists/{autoReceptionistId}/policies/{policyType}
func (p *PhoneAutoReceptionistsService) AddPolicySubSetting(ctx context.Context, pathParams *AddPolicySubSettingPathParams, body *AutoReceptionistPolicySubSetting) (*AutoReceptionistPolicySubSetting, *http.Response, error) {
	out := &AutoReceptionistPolicySubSetting{}

	res, err := p.client.request(ctx, http.MethodPost, fmt.Sprintf("/phone/auto_receptionists/%s/policies/%s", pathParams.autoReceptionistId, pathParams.policyType), nil, body, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type AddAutoReceptionistRequest struct {
	Name   string `json:"name"`
	SiteID string `json:"site_id,omitempty"`
}

type AddAutoReceptionistResponse struct {
	ExtensionNumber int    `json:"extension_number"`
	ID              string `json:"id"`
	Name            string `json:"name"`
}

// https://developers.zoom.us/docs/api/phone/#tag/auto-receptionists/post/phone/auto_receptionists
func (p *PhoneAutoReceptionistsService) AddAutoReceptionist(ctx context.Context, req *AddAutoReceptionistRequest) (*AddAutoReceptionistResponse, *http.Response, error) {
	out := &AddAutoReceptionistResponse{}

	res, err := p.client.request(ctx, http.MethodPost, "/phone/auto_receptionists", nil, req, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type AssignPhoneNumbersPathParams struct {
	autoReceptionistId string
}

type AssignPhoneNumbersRequest struct {
	PhoneNumbers []struct {
		ID     string `json:"id"`
		Number string `json:"number"`
	} `json:"phone_numbers"`
}

// https://developers.zoom.us/docs/api/phone/#tag/auto-receptionists/post/phone/auto_receptionists/{autoReceptionistId}/phone_numbers
func (p *PhoneAutoReceptionistsService) AssignPhoneNumbers(ctx context.Context, pathParams *AssignPhoneNumbersPathParams, body *AssignPhoneNumbersRequest) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodPost, fmt.Sprintf("/phone/auto_receptionists/%s/phone_numbers", pathParams.autoReceptionistId), nil, body, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type DeleteAutoReceptionistPathParams struct {
	autoReceptionistId string
}

// https://developers.zoom.us/docs/api/phone/#tag/auto-receptionists/delete/phone/auto_receptionists/{autoReceptionistId}
func (p *PhoneAutoReceptionistsService) DeleteAutoReceptionist(ctx context.Context, pathParams *DeleteAutoReceptionistPathParams) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodDelete, fmt.Sprintf("/phone/auto_receptionists/%s", pathParams.autoReceptionistId), nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type DeletePolicySubSettingPathParams struct {
	autoReceptionistId string
	policyType         string
}

type DeletePolicySubSettingQuery struct {
	SharedIDs []string `url:"shared_ids"`
}

// https://developers.zoom.us/docs/api/phone/#tag/auto-receptionists/delete/phone/auto_receptionists/{autoReceptionistId}/policies/{policyType}
func (p *PhoneAutoReceptionistsService) DeletePolicySubSetting(ctx context.Context, pathParams *DeletePolicySubSettingPathParams, query *DeletePolicySubSettingQuery) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodDelete, fmt.Sprintf("/phone/auto_receptionists/%s/policies/%s", pathParams.autoReceptionistId, pathParams.policyType), query, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type GetAutoReceptionistPathParams struct {
	autoReceptionistId string
}

type AutoReceptionistDetails struct {
	AudioPromptLanguage string `json:"audio_prompt_language"`
	CostCenter          string `json:"cost_center"`
	Department          string `json:"department"`
	ExtensionID         string `json:"extension_id"`
	ExtensionNumber     int    `json:"extension_number"`
	HolidayHours        []struct {
		From string `json:"from"`
		ID   string `json:"id"`
		Name string `json:"name"`
		To   string `json:"to"`
	} `json:"holiday_hours"`
	Name           string `json:"name"`
	OwnStorageName string `json:"own_storage_name"`
	PhoneNumbers   []struct {
		ID     string `json:"id"`
		Number string `json:"number"`
	} `json:"phone_numbers"`
	RecordingStorageLocation string `json:"recording_storage_location"`
	Site                     struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"site"`
	TimeZone string `json:"time_zone"`
}

// https://developers.zoom.us/docs/api/phone/#tag/auto-receptionists/get/phone/auto_receptionists/{autoReceptionistId}
func (p *PhoneAutoReceptionistsService) GetAutoReceptionist(ctx context.Context, pathParams *GetAutoReceptionistPathParams) (*AutoReceptionistDetails, *http.Response, error) {
	out := &AutoReceptionistDetails{}

	res, err := p.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/auto_receptionists/%s", pathParams.autoReceptionistId), nil, nil, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type GetAutoReceptionistPolicyPathParams struct {
	autoReceptionistId string
}

type AutoReceptionistPolicy struct {
	SMS struct {
		*AccountSettingStates
		InternationalSMS bool `json:"international_sms"`
	}
	VoicemailAccessMembers []struct {
		AccessUserID   string `json:"access_user_id"`
		AccessUserType string `json:"access_user_type"`
		Delete         bool   `json:"delete"`
		Download       bool   `json:"download"`
		SharedID       string `json:"shared_id"`
	} `json:"voicemail_access_members"`
	VoicemailNotificationByEmail struct {
		*AccountSettingStates
		Modified                      string `json:"modified"`
		ForwardVoicemailToEmail       bool   `json:"forward_voicemail_to_email"`
		IncludeVoicemailFile          bool   `json:"include_voicemail_file"`
		IncludeVoicemailTranscription bool   `json:"include_voicemail_transcription"`
	} `json:"voicemail_notification_by_email"`
	VoicemailTranscription struct {
		*AccountSettingStates
		Modified string `json:"modified"`
	} `json:"voicemail_transcription"`
}

// https://developers.zoom.us/docs/api/phone/#tag/auto-receptionists/get/phone/auto_receptionists/{autoReceptionistId}/policies
func (p *PhoneAutoReceptionistsService) GetAutoReceptionistPolicy(ctx context.Context, pathParams *GetAutoReceptionistPolicyPathParams) (*AutoReceptionistPolicy, *http.Response, error) {
	out := &AutoReceptionistPolicy{}

	res, err := p.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/auto_receptionists/%s/policies", pathParams.autoReceptionistId), nil, nil, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type ListAutoReceptionistsRequest struct {
	*PaginationOptions
}

type ListAutoReceptionistsResponse struct {
	*PaginationResponse
	AutoReceptionists []*AutoReceptionistDetails `json:"auto_receptionists"`
}

// https://developers.zoom.us/docs/api/phone/#tag/auto-receptionists/get/phone/auto_receptionists
func (p *PhoneAutoReceptionistsService) ListAutoReceptionists(ctx context.Context, req *ListAutoReceptionistsRequest) (*ListAutoReceptionistsResponse, *http.Response, error) {
	out := &ListAutoReceptionistsResponse{}

	res, err := p.client.request(ctx, http.MethodGet, "/phone/auto_receptionists", req, nil, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type UnassignPhoneNumbersPathParams struct {
	autoReceptionistId string
	phoneNumberId      string
}

// https://developers.zoom.us/docs/api/phone/#tag/auto-receptionists/delete/phone/auto_receptionists/{autoReceptionistId}/phone_numbers/{phoneNumberId}
func (p *PhoneAutoReceptionistsService) UnassignPhoneNumbers(ctx context.Context, pathParams *UnassignPhoneNumbersPathParams) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodDelete, fmt.Sprintf("/phone/auto_receptionists/%s/phone_numbers/%s", pathParams.autoReceptionistId, pathParams.phoneNumberId), nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type UnassignAllPhoneNumbersPathParams struct {
	autoReceptionistId string
}

// https://developers.zoom.us/docs/api/phone/#tag/auto-receptionists/delete/phone/auto_receptionists/{autoReceptionistId}/phone_numbers
func (p *PhoneAutoReceptionistsService) UnassignAllPhoneNumbers(ctx context.Context, pathParams *UnassignAllPhoneNumbersPathParams) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodDelete, fmt.Sprintf("/phone/auto_receptionists/%s/phone_numbers", pathParams.autoReceptionistId), nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type UpdatePolicySubSettingPathParams struct {
	autoReceptionistId string
	policyType         string
}

type UpdatePolicySubSettingRequest struct {
	VoiceMailAccessMember struct {
		AccessUserID   string `json:"access_user_id"`
		AccessUserType string `json:"access_user_type"`
		Delete         bool   `json:"delete"`
		Download       bool   `json:"download"`
		SharedID       string `json:"shared_id,omitempty"`
	} `json:"voice_mail_access_member"`
}

// https://developers.zoom.us/docs/api/phone/#tag/auto-receptionists/patch/phone/auto_receptionists/{autoReceptionistId}/policies/{policyType}
func (p *PhoneAutoReceptionistsService) UpdatePolicySubSetting(ctx context.Context, pathParams *UpdatePolicySubSettingPathParams, body *UpdatePolicySubSettingRequest) (*AutoReceptionistPolicySubSetting, *http.Response, error) {
	out := &AutoReceptionistPolicySubSetting{}

	res, err := p.client.request(ctx, http.MethodPatch, fmt.Sprintf("/phone/auto_receptionists/%s/policies/%s", pathParams.autoReceptionistId, pathParams.policyType), nil, body, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type UpdateAutoReceptionistPathParams struct {
	autoReceptionistId string
}

type UpdateAutoReceptionistRequest struct {
	AudioPromptLanguage      string `json:"audio_prompt_language,omitempty"`
	CostCenter               string `json:"cost_center,omitempty"`
	Department               string `json:"department,omitempty"`
	ExtensionNumber          int    `json:"extension_number,omitempty"`
	Name                     string `json:"name,omitempty"`
	RecordingStorageLocation string `json:"recording_storage_location,omitempty"`
	TimeZone                 string `json:"time_zone,omitempty"`
}

// https://developers.zoom.us/docs/api/phone/#tag/auto-receptionists/patch/phone/auto_receptionists/{autoReceptionistId}
func (p *PhoneAutoReceptionistsService) UpdateAutoReceptionist(ctx context.Context, pathParams *UpdateAutoReceptionistPathParams, body *UpdateAutoReceptionistRequest) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodPatch, fmt.Sprintf("/phone/auto_receptionists/%s", pathParams.autoReceptionistId), nil, body, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type UpdateAutoReceptionistPolicyPathParams struct {
	autoReceptionistId string
}

type UpdateAutoReceptionistPolicyRequest struct {
	SMS struct {
		Enable                    bool     `json:"enable"`
		InternationalSMS          bool     `json:"international_sms"`
		InternationalSMSCountries []string `json:"international_sms_countries,omitempty"`
		Reset                     bool     `json:"reset"`
	} `json:"sms,omitempty"`
	VoicemailNotificationByEmail struct {
		ForwardVoicemailToEmail       bool `json:"forward_voicemail_to_email"`
		IncludeVoicemailFile          bool `json:"include_voicemail_file"`
		IncludeVoicemailTranscription bool `json:"include_voicemail_transcription"`
	} `json:"voicemail_notification_by_email,omitempty"`
	VoicemailTranscription struct {
		Enable bool `json:"enable"`
		Reset  bool `json:"reset"`
	} `json:"voicemail_transcription,omitempty"`
}

// https://developers.zoom.us/docs/api/phone/#tag/auto-receptionists/patch/phone/auto_receptionists/{autoReceptionistId}/policies
func (p *PhoneAutoReceptionistsService) UpdateAutoReceptionistPolicy(ctx context.Context, pathParams *UpdateAutoReceptionistPolicyPathParams, body *UpdateAutoReceptionistPolicyRequest) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodPatch, fmt.Sprintf("/phone/auto_receptionists/%s/policies", pathParams.autoReceptionistId), nil, body, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type PhoneBillingAccountService struct {
	client *Client
}

type GetBillingAccountPathParams struct {
	billingAccountId string
}

type BillingAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// https://developers.zoom.us/docs/api/phone/#tag/billing-accounts/get/phone/billing_accounts/{billingAccountId}
func (p *PhoneBillingAccountService) GetBillingAccount(ctx context.Context, pathParams *GetBillingAccountPathParams) (*BillingAccount, *http.Response, error) {
	out := &BillingAccount{}

	res, err := p.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/billing_accounts/%s", pathParams.billingAccountId), nil, nil, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type ListBillingAccountsRequest struct {
	SiteID string `url:"site_id,omitempty"`
}

type ListBillingAccountsResponse struct {
	BillingAccounts []*BillingAccount `json:"billing_accounts"`
}

// https://developers.zoom.us/docs/api/phone/#tag/billing-accounts/get/phone/billing_accounts
func (p *PhoneBillingAccountService) ListBillingAccounts(ctx context.Context, req *ListBillingAccountsRequest) (*ListBillingAccountsResponse, *http.Response, error) {
	out := &ListBillingAccountsResponse{}

	res, err := p.client.request(ctx, http.MethodGet, "/phone/billing_accounts", req, nil, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type PhoneBlockedListService struct {
	client *Client
}

type CreateBlockedListRequest struct {
	BlockType   string `json:"block_type"`
	Comment     string `json:"comment,omitempty"`
	Country     string `json:"country,omitempty"`
	MatchType   string `json:"match_type,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Status      string `json:"status,omitempty"`
}

type CreateBlockedListResponse struct {
	ID string `json:"id"`
}

// https://developers.zoom.us/docs/api/phone/#tag/blocked-list/post/phone/blocked_list
func (p *PhoneBlockedListService) CreateBlockedList(ctx context.Context, req *CreateBlockedListRequest) (*CreateBlockedListResponse, *http.Response, error) {
	out := &CreateBlockedListResponse{}

	res, err := p.client.request(ctx, http.MethodPost, "/phone/blocked_list", nil, req, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type DeleteBlockedListPathParams struct {
	BlockedListId string
}

// https://developers.zoom.us/docs/api/phone/#tag/blocked-list/delete/phone/blocked_list/{blockedListId}
func (p *PhoneBlockedListService) DeleteBlockedList(ctx context.Context, pathParams *DeleteBlockedListPathParams) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodDelete, fmt.Sprintf("/phone/blocked_list/%s", pathParams.BlockedListId), nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type GetBlockedListPathParams struct {
	BlockedListId string
}

type BlockedListDetails struct {
	BlockType   string `json:"block_type"`
	Comment     string `json:"comment"`
	ID          string `json:"id"`
	MatchType   string `json:"match_type"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
}

// https://developers.zoom.us/docs/api/phone/#tag/blocked-list/get/phone/blocked_list/{blockedListId}
func (p *PhoneBlockedListService) GetBlockedList(ctx context.Context, pathParams *GetBlockedListPathParams) (*BlockedListDetails, *http.Response, error) {
	out := &BlockedListDetails{}

	res, err := p.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/blocked_list/%s", pathParams.BlockedListId), nil, nil, out)

	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type ListBlockedListRequest struct {
	*PaginationOptions
}

type ListBlockedListResponse struct {
	*PaginationResponse
	BlockedList []BlockedListDetails `json:"blocked_list"`
}

// https://developers.zoom.us/docs/api/phone/#tag/blocked-list/get/phone/blocked_list
func (p *PhoneBlockedListService) ListBlockedList(ctx context.Context, req *ListBlockedListRequest) (*ListBlockedListResponse, *http.Response, error) {
	out := &ListBlockedListResponse{}

	res, err := p.client.request(ctx, http.MethodGet, "/phone/blocked_list", req, nil, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type UpdateBlockedListPathParams struct {
	BlockedListId string
}

// https://developers.zoom.us/docs/api/phone/#tag/blocked-list/patch/phone/blocked_list/{blockedListId}
func (p *PhoneBlockedListService) UpdateBlockedList(ctx context.Context, pathParams *UpdateBlockedListPathParams, body *CreateBlockedListRequest) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodPatch, fmt.Sprintf("/phone/blocked_list/%s", pathParams.BlockedListId), nil, body, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type PhoneCallHandlingService struct {
	client *Client
}

type AddCallHandlingPathParams struct {
	ExtensionId string
	SettingType string
}

type AddCallHandlingRequest struct {
	Settings       interface{} `json:"settings"`
	SubSettingType string      `json:"sub_setting_type,omitempty"`
}

type CallForwardingSettings struct {
	Description string `json:"description,omitempty"`
	HolidayID   string `json:"holiday_id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type HolidaySettings struct {
	From string `json:"from,omitempty"`
	Name string `json:"name,omitempty"`
	To   string `json:"to,omitempty"`
}

type AddCallHandlingResponse struct {
	CallForwardingID string `json:"call_forwarding_id"`
	HolidayID        string `json:"holiday_id"`
}

// https://developers.zoom.us/docs/api/phone/#tag/call-handling/post/phone/extensions/{extensionId}/call_handling/{settingType}
func (p *PhoneCallHandlingService) AddCallHandling(ctx context.Context, pathParams *AddCallHandlingPathParams, req *AddCallHandlingRequest) (*AddCallHandlingResponse, *http.Response, error) {
	out := &AddCallHandlingResponse{}

	res, err := p.client.request(ctx, http.MethodPost, fmt.Sprintf("/phone/extensions/%s/call_handling/%s", pathParams.ExtensionId, pathParams.SettingType), nil, req, out)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return out, res, nil
}

type DeleteCallHandlingPathParams struct {
	ExtensionId string
	SettingType string
}

type DeleteCallHandlingQuery struct {
	CallForwardingId string `url:"call_forwarding_id,omitempty"`
	HolidayId        string `url:"holiday_id,omitempty"`
}

// https://developers.zoom.us/docs/api/phone/#tag/call-handling/delete/phone/extensions/{extensionId}/call_handling/{settingType}
func (p *PhoneCallHandlingService) DeleteCallHandling(ctx context.Context, pathParams *DeleteCallHandlingPathParams, query *DeleteCallHandlingQuery) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodDelete, fmt.Sprintf("/phone/extensions/%s/call_handling/%s", pathParams.ExtensionId, pathParams.SettingType), query, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}

	return res, nil
}

type GetCallHandlingSettingsPathParams struct {
	ExtensionId string
}

type BusinessHoursSettings struct {
	Settings struct {
		AllowMembersToReset  bool `json:"allow_members_to_reset"`
		AudioWhileConnecting struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"audio_while_connecting"`
		BusyRouting struct {
			Action                          int  `json:"action"`
			AllowCallersCheckVoicemail      bool `json:"allow_callers_check_voicemail"`
			BusyPlayCalleeVoicemailGreeting bool `json:"busy_play_callee_voicemail_greeting"`
			ConnectToOperator               bool `json:"connect_to_operator"`
			ForwardTo                       struct {
				Description       string `json:"description,omitempty"`
				DisplayName       string `json:"display_name,omitempty"`
				ExtensionID       string `json:"extension_id,omitempty"`
				ExtensionNumber   int    `json:"extension_number,omitempty"`
				ExtensionType     string `json:"extension_type,omitempty"`
				ID                string `json:"id,omitempty"`
				PhoneNumber       string `json:"phone_number,omitempty"`
				VoicemailGreeting struct {
					ID   string `json:"id,omitempty"`
					Name string `json:"name,omitempty"`
				} `json:"voicemail_greeting,omitempty"`
				MessageGreeting struct {
					ID   string `json:"id,omitempty"`
					Name string `json:"name,omitempty"`
				} `json:"message_greeting,omitempty"`
				Operator struct {
					DisplayName     string `json:"display_name,omitempty"`
					ExtensionID     string `json:"extension_id,omitempty"`
					ExtensionNumber int    `json:"extension_number,omitempty"`
					ExtensionType   string `json:"extension_type,omitempty"`
					ID              string `json:"id,omitempty"`
				} `json:"operator,omitempty"`
				PlayCalleeVoicemailGreeting   bool `json:"play_callee_voicemail_greeting,omitempty"`
				RequirePress1BeforeConnecting bool `json:"require_press_1_before_connecting,omitempty"`
				VoicemailLeavingInstructions  struct {
					ID   string `json:"id,omitempty"`
					Name string `json:"name,omitempty"`
				} `json:"voicemail_leaving_instructions,omitempty"`
				CallDistribution struct {
					HandleMultipleCalls           bool   `json:"handle_multiple_calls,omitempty"`
					RingDuration                  int    `json:"ring_duration,omitempty"`
					RingMode                      string `json:"ring_mode,omitempty"`
					SkipOfflineDevicePhoneNumbers bool   `json:"skip_offline_device_phone_numbers,omitempty"`
				} `json:"call_distribution,omitempty"`
			} `json:"forward_to"`
		} `json:"busy_routing"`
	} `json:"settings"`
	SubSettingType string `json:"sub_setting_type"`
}

type CallHandlingSettings struct {
}
