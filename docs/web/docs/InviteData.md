# InviteData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | invite id | 
**Type** | **string** |  | 
**Attributes** | [**InviteAttributes**](InviteAttributes.md) |  | 

## Methods

### NewInviteData

`func NewInviteData(id string, type_ string, attributes InviteAttributes, ) *InviteData`

NewInviteData instantiates a new InviteData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInviteDataWithDefaults

`func NewInviteDataWithDefaults() *InviteData`

NewInviteDataWithDefaults instantiates a new InviteData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *InviteData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *InviteData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *InviteData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *InviteData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *InviteData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *InviteData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *InviteData) GetAttributes() InviteAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *InviteData) GetAttributesOk() (*InviteAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *InviteData) SetAttributes(v InviteAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


