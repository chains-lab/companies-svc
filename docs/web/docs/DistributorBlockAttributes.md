# companyBlockAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**companyId** | [**uuid.UUID**](uuid.UUID.md) | ID of the company being blocked | 
**InitiatorId** | [**uuid.UUID**](uuid.UUID.md) | ID of the user who initiated the block | 
**Reason** | **string** | Reason for blocking the company | 
**Status** | **string** | Current status of the block | 
**BlockedAt** | **time.Time** | Timestamp when the block was initiated | 
**CancelledAt** | Pointer to **time.Time** | Timestamp when the block was lifted, if applicable | [optional] 

## Methods

### NewcompanyBlockAttributes

`func NewcompanyBlockAttributes(companyId uuid.UUID, initiatorId uuid.UUID, reason string, status string, blockedAt time.Time, ) *companyBlockAttributes`

NewcompanyBlockAttributes instantiates a new companyBlockAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewcompanyBlockAttributesWithDefaults

`func NewcompanyBlockAttributesWithDefaults() *companyBlockAttributes`

NewcompanyBlockAttributesWithDefaults instantiates a new companyBlockAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetcompanyId

`func (o *companyBlockAttributes) GetcompanyId() uuid.UUID`

GetcompanyId returns the companyId field if non-nil, zero value otherwise.

### GetcompanyIdOk

`func (o *companyBlockAttributes) GetcompanyIdOk() (*uuid.UUID, bool)`

GetcompanyIdOk returns a tuple with the companyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetcompanyId

`func (o *companyBlockAttributes) SetcompanyId(v uuid.UUID)`

SetcompanyId sets companyId field to given value.


### GetInitiatorId

`func (o *companyBlockAttributes) GetInitiatorId() uuid.UUID`

GetInitiatorId returns the InitiatorId field if non-nil, zero value otherwise.

### GetInitiatorIdOk

`func (o *companyBlockAttributes) GetInitiatorIdOk() (*uuid.UUID, bool)`

GetInitiatorIdOk returns a tuple with the InitiatorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitiatorId

`func (o *companyBlockAttributes) SetInitiatorId(v uuid.UUID)`

SetInitiatorId sets InitiatorId field to given value.


### GetReason

`func (o *companyBlockAttributes) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *companyBlockAttributes) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *companyBlockAttributes) SetReason(v string)`

SetReason sets Reason field to given value.


### GetStatus

`func (o *companyBlockAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *companyBlockAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *companyBlockAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetBlockedAt

`func (o *companyBlockAttributes) GetBlockedAt() time.Time`

GetBlockedAt returns the BlockedAt field if non-nil, zero value otherwise.

### GetBlockedAtOk

`func (o *companyBlockAttributes) GetBlockedAtOk() (*time.Time, bool)`

GetBlockedAtOk returns a tuple with the BlockedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockedAt

`func (o *companyBlockAttributes) SetBlockedAt(v time.Time)`

SetBlockedAt sets BlockedAt field to given value.


### GetCancelledAt

`func (o *companyBlockAttributes) GetCancelledAt() time.Time`

GetCancelledAt returns the CancelledAt field if non-nil, zero value otherwise.

### GetCancelledAtOk

`func (o *companyBlockAttributes) GetCancelledAtOk() (*time.Time, bool)`

GetCancelledAtOk returns a tuple with the CancelledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCancelledAt

`func (o *companyBlockAttributes) SetCancelledAt(v time.Time)`

SetCancelledAt sets CancelledAt field to given value.

### HasCancelledAt

`func (o *companyBlockAttributes) HasCancelledAt() bool`

HasCancelledAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


