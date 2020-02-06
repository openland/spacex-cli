/* tslint:disable */
/* eslint-disable */
type Maybe<T> = T | null;
type MaybeInput<T> = T | null | undefined;
type Inline<V> = {} | V;

// Enums
export enum DeviceKind {
    LAMP = 'LAMP',
    SWITCH = 'SWITCH',
    UNKNOWN = 'UNKNOWN',
}

// Input Types
export interface DeviceDescription {
    name?: MaybeInput<string>;
    description?: MaybeInput<string>;
    icon?: MaybeInput<ImageRef>;
}
export interface ImageRef {
    url: string;
}

// Fragments
export type DeviceNano = (
    & { __typename: 'Lamp' | 'Switch' | 'Lock' | string }
    & { id: string}
);
export type LampShort = (
    & { __typename: 'Lamp' }
    & { id: string}
    & { brightness: number}
    & { minBrightness: number}
    & { maxBrightness: number}
);
export type DeviceShort = (
    & { __typename: 'Lamp' | 'Switch' | 'Lock' | string }
    & { description: Maybe<string>}
    & { addedBy: (
        & { __typename: 'User' }
        & { id: string}
        & { username: string}
    )}
    & DeviceNano
    & Inline<(
        & { __typename: 'Lamp' }
        & LampShort
    )>
    & Inline<(
        & { __typename: 'Switch' }
        & { id: string}
        & { numberOfButtons: Maybe<number>}
    )>
    & Inline<(
        & { __typename: 'Lamp' }
        & { id: string}
        & { brightness: number}
    )>
);
export type UserShort = (
    & { __typename: 'User' }
    & { id: string}
);

// Queries
export type DiscoverDevices = (
    & { discover: [(
        & { __typename: 'DiscoveredThing' }
        & { id: string}
        & { name: string}
        & { host: string}
        & { port: number}
    )]}
);

// Mutations
export interface AddUserVariables {
    username: string;
    password: string;
}
export type AddUser = (
    & { addUser: (
        & { __typename: 'User' }
        & UserShort
    )}
);
export interface UpdateDeviceDescriptionVariables {
    id: string;
    description: DeviceDescription;
}
export type UpdateDeviceDescription = (
    & { updateDeviceDescription: (
        & { __typename: 'Lamp' | 'Switch' | 'Lock' | string }
        & DeviceShort
    )}
);

// Subscriptions