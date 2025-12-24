package zoom

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

type PhoneService struct {
	client   *Client
	Accounts *PhoneAccountsService
	Alerts   *PhoneAlertsService
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
	RuleConditions    struct {
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
}
