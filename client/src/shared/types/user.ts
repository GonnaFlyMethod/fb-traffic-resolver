export interface TUser {
  email: string;
  password: string;
}

export interface TUserMeta extends TUser {
  id?: number;
  name: string;
  surname: string;
  confirmPassword?: string;
  avatar?: string;
}
