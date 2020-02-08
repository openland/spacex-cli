export type Maybe<T> = T | null;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string,
  String: string,
  Boolean: boolean,
  Int: number,
  Float: number,
};

export type Device = {
  id: Scalars['ID'],
  description?: Maybe<Scalars['String']>,
  addedBy: User,
};

export type DeviceDescription = {
  name?: Maybe<Scalars['String']>,
  description?: Maybe<Scalars['String']>,
  icon?: Maybe<ImageRef>,
};

export enum DeviceKind {
  Lamp = 'LAMP',
  Switch = 'SWITCH',
  Unknown = 'UNKNOWN'
}

export type DiscoveredThing = {
   __typename?: 'DiscoveredThing',
  id: Scalars['ID'],
  name: Scalars['String'],
  host: Scalars['String'],
  port: Scalars['Int'],
};

export type ImageRef = {
  url: Scalars['String'],
};

export type Lamp = Device & {
   __typename?: 'Lamp',
  id: Scalars['ID'],
  description?: Maybe<Scalars['String']>,
  addedBy: User,
  brightness: Scalars['Float'],
  minBrightness: Scalars['Float'],
  maxBrightness: Scalars['Float'],
};

export type Lock = Device & {
   __typename?: 'Lock',
  id: Scalars['ID'],
  description?: Maybe<Scalars['String']>,
  addedBy: User,
  locked: Scalars['Boolean'],
};

export type Mutation = {
   __typename?: 'Mutation',
  addUser: User,
  updateDeviceDescription: Device,
};


export type MutationAddUserArgs = {
  username: Scalars['String'],
  password: Scalars['String']
};


export type MutationUpdateDeviceDescriptionArgs = {
  id: Scalars['ID'],
  description: DeviceDescription
};

export type Query = {
   __typename?: 'Query',
  discover: Array<DiscoveredThing>,
  users: Array<User>,
  allDevices: Array<Device>,
};

export type Switch = Device & {
   __typename?: 'Switch',
  id: Scalars['ID'],
  description?: Maybe<Scalars['String']>,
  addedBy: User,
  numberOfButtons?: Maybe<Scalars['Int']>,
};

export type User = {
   __typename?: 'User',
  id: Scalars['ID'],
  username: Scalars['String'],
};

export type DiscoverDevicesQueryVariables = {};


export type DiscoverDevicesQuery = (
  { __typename?: 'Query' }
  & { discover: Array<(
    { __typename?: 'DiscoveredThing' }
    & Pick<DiscoveredThing, 'id' | 'name' | 'host' | 'port'>
  )> }
);

export type UpdateDeviceDescriptionMutationVariables = {
  id: Scalars['ID'],
  description: DeviceDescription
};


export type UpdateDeviceDescriptionMutation = (
  { __typename?: 'Mutation' }
  & { updateDeviceDescription: (
    { __typename?: 'Lamp' }
    & DeviceShort_Lamp_Fragment
  ) | (
    { __typename?: 'Switch' }
    & DeviceShort_Switch_Fragment
  ) | (
    { __typename?: 'Lock' }
    & DeviceShort_Lock_Fragment
  ) }
);

export type AddUserMutationVariables = {
  username: Scalars['String'],
  password: Scalars['String']
};


export type AddUserMutation = (
  { __typename?: 'Mutation' }
  & { addUser: (
    { __typename?: 'User' }
    & UserShortFragment
  ) }
);

export type UserShortFragment = (
  { __typename?: 'User' }
  & Pick<User, 'id'>
);

export type LampShortFragment = (
  { __typename?: 'Lamp' }
  & Pick<Lamp, 'id' | 'brightness' | 'minBrightness' | 'maxBrightness'>
);

type DeviceNano_Lamp_Fragment = (
  { __typename?: 'Lamp' }
  & Pick<Lamp, 'id'>
);

type DeviceNano_Switch_Fragment = (
  { __typename?: 'Switch' }
  & Pick<Switch, 'id'>
);

type DeviceNano_Lock_Fragment = (
  { __typename?: 'Lock' }
  & Pick<Lock, 'id'>
);

export type DeviceNanoFragment = DeviceNano_Lamp_Fragment | DeviceNano_Switch_Fragment | DeviceNano_Lock_Fragment;

type DeviceShort_Lamp_Fragment = (
  { __typename?: 'Lamp' }
  & Pick<Lamp, 'id' | 'brightness' | 'description'>
  & { addedBy: (
    { __typename: 'User' }
    & Pick<User, 'id' | 'username'>
  ) }
  & LampShortFragment
  & DeviceNano_Lamp_Fragment
);

type DeviceShort_Switch_Fragment = (
  { __typename?: 'Switch' }
  & Pick<Switch, 'id' | 'numberOfButtons' | 'description'>
  & { addedBy: (
    { __typename: 'User' }
    & Pick<User, 'id' | 'username'>
  ) }
  & DeviceNano_Switch_Fragment
);

type DeviceShort_Lock_Fragment = (
  { __typename?: 'Lock' }
  & Pick<Lock, 'description'>
  & { addedBy: (
    { __typename: 'User' }
    & Pick<User, 'id' | 'username'>
  ) }
  & DeviceNano_Lock_Fragment
);

export type DeviceShortFragment = DeviceShort_Lamp_Fragment | DeviceShort_Switch_Fragment | DeviceShort_Lock_Fragment;
