# InitDistributorBlockDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DistributorId** | **string** | The UUID of the distributor to be blocked. | 
**Reason** | **string** | The reason for blocking the distributor. | 

## Methods

### NewInitDistributorBlockDataAttributes

`func NewInitDistributorBlockDataAttributes(distributorId string, reason string, ) *InitDistributorBlockDataAttributes`

NewInitDistributorBlockDataAttributes instantiates a new InitDistributorBlockDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInitDistributorBlockDataAttributesWithDefaults

`func NewInitDistributorBlockDataAttributesWithDefaults() *InitDistributorBlockDataAttributes`

NewInitDistributorBlockDataAttributesWithDefaults instantiates a new InitDistributorBlockDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDistributorId

`func (o *InitDistributorBlockDataAttributes) GetDistributorId() string`

GetDistributorId returns the DistributorId field if non-nil, zero value otherwise.

### GetDistributorIdOk

`func (o *InitDistributorBlockDataAttributes) GetDistributorIdOk() (*string, bool)`

GetDistributorIdOk returns a tuple with the DistributorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistributorId

`func (o *InitDistributorBlockDataAttributes) SetDistributorId(v string)`

SetDistributorId sets DistributorId field to given value.


### GetReason

`func (o *InitDistributorBlockDataAttributes) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *InitDistributorBlockDataAttributes) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *InitDistributorBlockDataAttributes) SetReason(v string)`

SetReason sets Reason field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


