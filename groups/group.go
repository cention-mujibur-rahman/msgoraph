package groups

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cention-mujibur-rahman/msgoraph/client"
	"github.com/cention-mujibur-rahman/msgoraph/internal"
)

// ServiceContext represents a namespace under which all of the operations against user-namespaced
// resources are accessed.
type ServiceContext struct {
	client client.Client
}

// Service creates a new users.ServiceContext with the given authentication credentials.
func Service(client client.Client) *ServiceContext {
	return &ServiceContext{client: client}
}

// CreateGroupRequest is all the available args you can set when creating a user.
type CreateGroupRequest struct {
	MailEnabled     bool   `json:"mailEnabled"`
	DisplayName     string `json:"displayName"`
	Description     string `json:"description"`
	MailNickname    string `json:"mailNickname"`
	Visibility      string `json:"visibility"`
	SecurityEnabled bool   `json:"securityEnabled"`
}

// GetGroupResponse is the response to expect on a GetGroup Request.
type GetGroupResponse struct {
	Context string `json:"@odata.context"`
	Group
}

// Group Create a new group as specified in the request body. You can create one of the following groups
// documentation https://docs.microsoft.com/en-us/graph/api/group-post-groups?view=graph-rest-beta&tabs=http
type Group struct {
	ID           *string `json:"id"`
	Mail         *string `json:"mail"`
	MailNickname *string `json:"mailNickname"`
}

// CreateGroup creates a new groups in the tenant.
func (s *ServiceContext) CreateGroup(createGroup CreateGroupRequest) (Group, error) {
	body, err := internal.GraphRequest(s.client, "POST", "v1.0/groups", nil, createGroup)
	if err != nil {
		log.Printf("Error GraphRequest %#v", err)
		return Group{}, err
	}
	var data GetGroupResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Group{}, err
	}
	return data.Group, nil
}

// CreateGroupsTeams Create a new team from a group.
//In order to create a team, the group must have a least one owner.
//If the group was created less than 15 minutes ago, it's possible for the Create team call to fail with a 404 error code due to replication delays.
//The recommended pattern is to retry the Create team call three times, with a 10 second delay between calls.
//
//https://docs.microsoft.com/en-us/graph/api/team-put-teams?view=graph-rest-beta&tabs=http
func (s *ServiceContext) CreateGroupsTeams(payloadBody interface{}) (Group, error) {
	body, err := internal.GraphRequest(s.client, "POST", "v1.0/groups", nil, payloadBody)
	if err != nil {
		log.Printf("Error CreateGroupsTeams GraphRequest %#v", err)
		return Group{}, err
	}
	var data GetGroupResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Group{}, err
	}
	return data.Group, nil
}

// GetGroupsTeams Get all groups.
func (s *ServiceContext) GetGroupsTeams() ([]Group, error) {
	body, err := internal.GraphRequest(s.client, "GET", "v1.0/groups", nil, nil)
	if err != nil {
		log.Printf("Error GetGroupsTeams GraphRequest %#v", err)
		return nil, err
	}
	var data GetAllGroupResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data.Value, nil
}

// GetAllGroupResponse is the response to expect on a GetGroup Request.
type GetAllGroupResponse struct {
	Context string `json:"@odata.context"`
	Value   []Group
}

// GetGroupsChannels Get all channels under a group.
//If successful, this method returns a 200 OK response code and collection of Channel objects in the response body.
//
//https://docs.microsoft.com/en-us/graph/api/channel-list?view=graph-rest-beta&tabs=http
func (s *ServiceContext) GetGroupsChannels(groupID string) ([]Channel, error) {
	url := fmt.Sprintf("v1.0/teams/%v/channels", groupID)
	body, err := internal.GraphRequest(s.client, "GET", url, nil, nil)
	if err != nil {
		log.Printf("Error GetGroupsChannels GraphRequest %#v", err)
		return nil, err
	}
	var data GetAllChannelResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data.Value, nil
}

// GetAllChannelResponse is the response to expect on a GetGroupsChannels Request.
type GetAllChannelResponse struct {
	Context string `json:"@odata.context"`
	Value   []Channel
}

//Channel represents a channel struct
//
//https://docs.microsoft.com/en-us/graph/api/channel-list?view=graph-rest-beta&tabs=http
type Channel struct {
	ID             *string `json:"id"`
	Description    *string `json:"description"`
	DisplayName    *string `json:"displayName"`
	WebURL         *string `json:"webUrl"`
	MembershipType *string `json:"membershipType"`
}

//GetChannelsContact List members
//If successful, this method returns a 200 OK response code and a collection of conversationMember objects in the response body.
//
//https://docs.microsoft.com/en-us/graph/api/team-list-members?view=graph-rest-beta&tabs=http
func (s *ServiceContext) GetChannelsContact(groupID string) ([]Contact, error) {
	url := fmt.Sprintf("v1.0/teams/%v/members", groupID)
	body, err := internal.GraphRequest(s.client, "GET", url, nil, nil)
	if err != nil {
		log.Printf("Error GetChannelsContact GraphRequest %#v", err)
		return nil, err
	}
	var data GetAllContactResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data.Value, nil
}

// GetAllContactResponse is the response to expect on a GetChannelsContact Request.
type GetAllContactResponse struct {
	Context string `json:"@odata.context"`
	Count   int    `json:"@odata.count"`
	Value   []Contact
}

//Contact represents a contact struct
type Contact struct {
	ID          *string `json:"id"`
	UserID      *string `json:"userId"`
	DisplayName *string `json:"displayName"`
	Email       *string `json:"email"`
}

// GetMessageResponse is the response to expect on a GetChannelsContact Request.
type GetMessageResponse struct {
	Context string `json:"@odata.context"`
	Count   int    `json:"@odata.count"`
	Value   []ChannelMessage
}

//ChannelMessage represents a contact struct
type ChannelMessage struct {
	ID              *string               `json:"id"`
	ReplyToID       *string               `json:"replyToId"`
	MessageType     *string               `json:"messageType"`
	CreatedDateTime *string               `json:"createdDateTime"`
	Subject         *string               `json:"subject"`
	Summary         *string               `json:"summary"`
	From            *From                 `json:"from"`
	Body            *MessageBody          `json:"body"`
	Attachments     []*MessageAttachments `json:"attachments"`
	Mentions        []*MessageMentions    `json:"mentions"`
	ChannelIdentity *ChannelIdentity      `json:"channelIdentity"`
}
type ChannelIdentity struct {
	TeamID    *string `json:"teamId"`
	ChannelID *string `json:"channelId"`
}
type MessageAttachments struct {
	ID           *string `json:"id"`
	ContentType  *string `json:"contentType"`
	ContentURL   *string `json:"contentUrl"`
	Content      *string `json:"content"`
	Name         *string `json:"name"`
	ThumbnailURL *string `json:"thumbnailUrl"`
}

type MessageMentions struct {
	ID          *int              `json:"id"`
	MentionText *string           `json:"mentionText"`
	Mentioned   *MentionedMessage `json:"mentioned"`
}

type MentionedMessage struct {
	User *User `json:"user"`
}

type MessageBody struct {
	ContentType *string `json:"contentType"`
	Content     *string `json:"content"`
}

type From struct {
	User *User `json:"user"`
}

type User struct {
	ID          *string `json:"id"`
	DisplayName *string `json:"displayName"`
}

//GetTeamsMessage List channel messages
//If successful, this method returns a 200 OK response code and a collection of chatMessage objects in the response body.
//
//https://docs.microsoft.com/en-us/graph/api/channel-list-messages?view=graph-rest-beta&tabs=http
func (s *ServiceContext) GetTeamsMessage(groupID, channelID string) (GetMessageResponse, error) {
	var data GetMessageResponse
	url := fmt.Sprintf("beta/teams/%v/channels/%v/messages", groupID, channelID)
	body, err := internal.GraphRequest(s.client, "GET", url, nil, nil)
	if err != nil {
		log.Printf("Error GetTeamsMessage GraphRequest %#v", err)
		return data, err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//GetTeamsMessageReplies Get a single reply to a message in a channel of a team.
//
//https://docs.microsoft.com/en-us/graph/api/channel-get-messagereply?view=graph-rest-beta&tabs=http
func (s *ServiceContext) GetTeamsMessageReplies(groupID, channelID, messageID string) (GetMessageResponse, error) {
	var data GetMessageResponse
	url := fmt.Sprintf("beta/teams/%v/channels/%v/messages/%v/replies", groupID, channelID, messageID)
	body, err := internal.GraphRequest(s.client, "GET", url, nil, nil)
	if err != nil {
		log.Printf("Error GetTeamsMessage GraphRequest %#v", err)
		return data, err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//SendTeamsMessage Sends channel messages
//If successful, this method returns a 200 OK response code and a collection of chatMessage objects in the response body.
//
//https://docs.microsoft.com/en-us/graph/api/channel-post-messages?view=graph-rest-beta&tabs=http
func (s *ServiceContext) SendTeamsMessage(groupID, channelID string, payloadBody interface{}) (ChannelMessage, error) {
	var data ChannelMessage
	url := fmt.Sprintf("beta/teams/%v/channels/%v/messages", groupID, channelID)
	body, err := internal.GraphRequest(s.client, "POST", url, nil, payloadBody)
	if err != nil {
		log.Printf("Error SendTeamsMessage GraphRequest %#v", err)
		return data, err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
