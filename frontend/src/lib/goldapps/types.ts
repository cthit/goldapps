// Server types
export interface Suggestions {
  userChanges?: UserActions;
  groupChanges?: GroupActions;
}

export interface UserActions {
  userUpdates?: UserUpdate[];
  additions?: User[];
  deletions?: User[];
}

export interface User {
  cid: string;
  first_name: string;
  second_name: string;
  nick: string;
  mail: string;
}

export interface UserUpdate {
  before: User;
  after: User;
}

export interface GroupActions {
  groupUpdates?: GroupUpdate[];
  additions?: Group[];
  deletions?: Group[];
}

export interface Group {
  email: string;
  type: string;
  members?: string[];
  aliases?: string[];
  expendable: boolean;
}

export interface GroupUpdate {
  before: Group;
  after: Group;
}

// Client types
export interface SuggestionAddUser {
  type: "AddUser";
  user: User;
}

export interface SuggestionDeleteUser {
  type: "DeleteUser";
  user: User;
}

export interface SuggestionChangeUser {
  type: "ChangeUser";
  before: User;
  after: User;
}

export type SuggestionUser =
  | SuggestionAddUser
  | SuggestionDeleteUser
  | SuggestionChangeUser;

export interface SuggestionAddGroup {
  type: "AddGroup";
  group: Group;
}

export interface SuggestionDeleteGroup {
  type: "DeleteGroup";
  group: Group;
}

export interface SuggestionChangeGroup {
  type: "ChangeGroup";
  before: Group;
  after: Group;
}

export type SuggestionGroup =
  | SuggestionAddGroup
  | SuggestionDeleteGroup
  | SuggestionChangeGroup;

export type Suggestion = SuggestionUser | SuggestionGroup;

export type SuggestionWithState = Suggestion & {
  selected: boolean;
  error: string | null;
};
