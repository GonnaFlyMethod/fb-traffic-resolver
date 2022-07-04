import { makeAutoObservable } from "mobx";

import { API } from "services";
import { TUserMeta } from "shared/types";
import LoadingModel from "models/Loading";

class UsersModel {
  private _users: TUserMeta[] = [];

  loading: LoadingModel;

  constructor() {
    makeAutoObservable(this);

    this.loading = new LoadingModel();
  }

  set users(data: TUserMeta[]) {
    this._users = data;
  }

  get users() {
    return this._users;
  }

  async fetch() {
    this.loading.begin();

    const data = await API.user.list();
    this.users = data.users;

    this.loading.end();
  }

  async create(user: TUserMeta) {
    this.loading.begin();

    await API.user.create(user);
    this.fetch();

    this.loading.end();
  }

  async remove(id: number) {
    this.loading.begin();

    await API.user.remove(id);
    this.fetch();

    this.loading.end();
  }

  async update(user: TUserMeta) {
    this.loading.begin();

    await API.user.update(user);
    this.fetch();

    this.loading.end();
  }
}

export default new UsersModel();
