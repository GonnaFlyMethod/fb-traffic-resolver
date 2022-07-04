import { createContext, useContext } from "react";

import LoadingModel from "models/Loading";
import { fetch } from "services";

export class RootStore {
  readonly loading: LoadingModel;

  constructor() {
    this.init();

    this.loading = new LoadingModel();
  }

  async init() {
    fetch.init();
  }
}

const rootStore = new RootStore();

const rootStoreContext = createContext<RootStore>(rootStore);

export const RootStoreProvider = ({ children }: { children: JSX.Element }) => (
  <rootStoreContext.Provider value={rootStore}>
    {children}
  </rootStoreContext.Provider>
);
export const useRootStore = () => useContext(rootStoreContext);

export default rootStore;
