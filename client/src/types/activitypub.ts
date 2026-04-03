// ActivityPub type constants
export const ACTIVITY_TYPES = [
  'Accept', 'Add', 'Announce', 'Arrive', 'Block', 'Create', 'Delete',
  'Dislike', 'Flag', 'Follow', 'Ignore', 'Invite', 'Join', 'Leave', 'Like',
  'Listen', 'Move', 'Offer', 'Question', 'Reject', 'Read', 'Remove',
  'TentativeReject', 'TentativeAccept', 'Travel', 'Undo', 'Update', 'View',
] as const

export const ACTOR_TYPES = ['Application', 'Group', 'Organization', 'Person', 'Service'] as const
export const OBJECT_TYPES = [
  'Article', 'Audio', 'Document', 'Event', 'Image', 'Note', 'Page',
  'Place', 'Profile', 'Relationship', 'Tombstone', 'Video',
] as const

export type ActivityType = (typeof ACTIVITY_TYPES)[number]
export type ActorType = (typeof ACTOR_TYPES)[number]
export type ObjectType = (typeof OBJECT_TYPES)[number]

// Base AP object
export interface APObject {
  '@context'?: string | string[] | Record<string, unknown>
  id: string
  type: string
  attachment?: APAttachment | APAttachment[]
  attributedTo?: string | string[] | APActor | APActor[]
  content?: string
  name?: string
  endTime?: string
  icon?: APImage
  image?: APImage | APImage[]
  inReplyTo?: string | APObject
  published?: string
  replies?: string | APCollection
  summary?: string
  tag?: APTag[]
  updated?: string
  url?: string | string[]
  to?: string[]
  cc?: string[]
  mediaType?: string
  duration?: string
  sensitive?: boolean
}

export interface APImage {
  type: 'Image'
  url: string
  mediaType?: string
  width?: number
  height?: number
}

export interface APAttachment {
  type: string
  url: string
  mediaType?: string
  name?: string
  width?: number
  height?: number
}

export interface APTag {
  type: string
  href?: string
  name?: string
}

export interface APPublicKey {
  id: string
  owner: string
  publicKeyPem: string
}

// Actor (Person, Group, etc.)
export interface APActor extends APObject {
  type: ActorType
  inbox: string
  outbox: string
  following?: string
  followers?: string
  liked?: string
  preferredUsername?: string
  manuallyApprovesFollowers?: boolean
  publicKey?: APPublicKey
  endpoints?: {
    sharedInbox?: string
    uploadMedia?: string
  }
}

// Activity
export interface APActivity extends APObject {
  type: ActivityType
  actor: string | APActor
  object?: string | APObject | APActivity | APActor
}

// Collections
export interface APCollection extends APObject {
  type: 'OrderedCollection' | 'Collection'
  totalItems: number
  first?: string | APCollectionPage
  last?: string
  orderedItems?: (string | APObject | APActivity | APActor)[]
  items?: (string | APObject | APActivity | APActor)[]
}

export interface APCollectionPage extends APObject {
  type: 'OrderedCollectionPage' | 'CollectionPage'
  partOf?: string
  orderedItems?: (string | APObject | APActivity | APActor)[]
  items?: (string | APObject | APActivity | APActor)[]
  next?: string
  prev?: string
  first?: string
  last?: string
}

// App-level auth user type
export interface AuthUser {
  username: string
}

// Check response from /check endpoint
export interface CheckResponse {
  exists: boolean
  iri: string
}

// Collection item type discriminator
export type CollectionItemType = 'Activity' | 'Object' | 'Actor'
