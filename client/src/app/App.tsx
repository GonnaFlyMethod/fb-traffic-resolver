import React from "react";
import { BrowserRouter } from "react-router-dom";
import { ModalProvider } from "react-modal-hook";
import { Slide, toast } from "react-toastify";
import { configure } from "mobx";
import { observer } from "mobx-react-lite";

import "react-toastify/dist/ReactToastify.css";

import { RootStoreProvider } from "stores";

import { Layout } from "./Layout";
import Router from "./Router";

configure({
  reactionScheduler: (f) => setTimeout(f),
});

toast.configure({
  position: "top-right",
  autoClose: 5000,
  hideProgressBar: false,
  closeOnClick: true,
  pauseOnHover: true,
  draggable: false,
  transition: Slide,
});

function App() {
  return (
    <React.StrictMode>
      <BrowserRouter>
        <RootStoreProvider>
          <ModalProvider>
            <Layout>
              <Router />
            </Layout>
          </ModalProvider>
        </RootStoreProvider>
      </BrowserRouter>
    </React.StrictMode>
  );
}

export default observer(App);
