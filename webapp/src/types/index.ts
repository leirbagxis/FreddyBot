export interface User {
  id: number;
  first_name: string;
  username: string;
  is_admin: boolean;
  channels: Channel[];
}

export interface Channel {
  id: number;
  title: string;
  newPackCaption: string;
  inviteUrl: string;
  ownerId: number;
  defaultCaption?: DefaultCaption;
  buttons: Button[];
}

export interface DefaultCaption {
  captionId: string;
  caption: string;
  messagePermission?: MessagePermission;
  buttonsPermission?: ButtonsPermission;
}

export interface Button {
  buttonId: string;
  nameButton: string;
  buttonUrl: string;
  positionX: number;
  positionY: number;
}

export interface MessagePermission {
  linkPreview: boolean;
  message: boolean;
  audio: boolean;
  video: boolean;
  photo: boolean;
  sticker: boolean;
  gif: boolean;
}

export interface ButtonsPermission {
  message: boolean;
  audio: boolean;
  video: boolean;
  photo: boolean;
  sticker: boolean;
  gif: boolean;
}
