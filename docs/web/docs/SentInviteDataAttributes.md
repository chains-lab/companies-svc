# SentInviteDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DistributorId** | [**uuid.UUID**](uuid.UUID.md) | ID of the distributor the invite is for | 
**Role** | **string** | Role assigned to the invited user | 

## Methods

### NewSentInviteDataAttributes

`func NewSentInviteDataAttributes(distributorId uuid.UUID, role string, ) *SentInviteDataAttributes`

NewSentInviteDataAttributes instantiates a new SentInviteDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSentInviteDataAttributesWithDefaults

`func NewSentInviteDataAttributesWithDefaults() *SentInviteDataAttributes`

NewSentInviteDataAttributesWithDefaults instantiates a new SentInviteDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDistributorId

`func (o *SentInviteDataAttributes) GetDistributorId() uuid.UUID`

GetDistributorId returns the DistributorId field if non-nil, zero value otherwise.

### GetDistributorIdOk

`func (o *SentInviteDataAttributes) GetDistributorIdOk() (*uuid.UUID, bool)`

GetDistributorIdOk returns a tuple with the DistributorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistributorId

`func (o *SentInviteDataAttributes) SetDistributorId(v uuid.UUID)`

SetDistributorId sets DistributorId field to given value.


### GetRole

`func (o *SentInviteDataAttributes) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *SentInviteDataAttributes) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *SentInviteDataAttributes) SetRole(v string)`

SetRole sets Role field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


