import { fetch } from "services";
import { TUserMeta } from "shared/types";

export const create = (user: TUserMeta) =>
  fetch.post<TUserMeta>("/user/create", user);

export const update = (user: TUserMeta) => fetch.post("/user/update", user);

export const remove = (id: number) => fetch.post("/user/remove", { id });

export const list = () =>
  fetch.post<{ users: TUserMeta[]; count: number }>("/user/list", {});
