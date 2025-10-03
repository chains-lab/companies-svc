# CreatecompanyBlockDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**companyId** | [**uuid.UUID**](uuid.UUID.md) | The UUID of the company to be blocked. | 
**Reason** | **string** | The reason for blocking the company. | 

## Methods

### NewCreatecompanyBlockDataAttributes

`func NewCreatecompanyBlockDataAttributes(companyId uuid.UUID, reason string, ) *CreatecompanyBlockDataAttributes`

NewCreatecompanyBlockDataAttributes instantiates a new CreatecompanyBlockDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreatecompanyBlockDataAttributesWithDefaults

`func NewCreatecompanyBlockDataAttributesWithDefaults() *CreatecompanyBlockDataAttributes`

NewCreatecompanyBlockDataAttributesWithDefaults instantiates a new CreatecompanyBlockDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetcompanyId

`func (o *CreatecompanyBlockDataAttributes) GetcompanyId() uuid.UUID`

GetcompanyId returns the companyId field if non-nil, zero value otherwise.

### GetcompanyIdOk

`func (o *CreatecompanyBlockDataAttributes) GetcompanyIdOk() (*uuid.UUID, bool)`

GetcompanyIdOk returns a tuple with the companyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetcompanyId

`func (o *CreatecompanyBlockDataAttributes) SetcompanyId(v uuid.UUID)`

SetcompanyId sets companyId field to given value.


### GetReason

`func (o *CreatecompanyBlockDataAttributes) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *CreatecompanyBlockDataAttributes) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *CreatecompanyBlockDataAttributes) SetReason(v string)`

SetReason sets Reason field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


