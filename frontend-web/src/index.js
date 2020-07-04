import "react-app-polyfill/ie11"; // For IE 11 support
import "react-app-polyfill/stable";
import "./polyfill";
import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import * as serviceWorker from "./serviceWorker";

import { icons } from "./assets/icons";

import { Provider } from "react-redux";
import store from "./store";

import { I18nextProvider } from "react-i18next";
import i18next from "i18next";

import assets_en from "./translations/en/assets.json";
import assets_id from "./translations/id/assets.json";

import navigation_en from "./translations/en/navigation.json";
import navigation_id from "./translations/id/navigation.json";

i18next.init({
  interpolation: { escapeValue: false },
  lng: "id",
  resources: {
    en: {
      assets: assets_en,
      navigation: navigation_en,
    },
    id: {
      assets: assets_id,
      navigation: navigation_id,
    },
  },
});

React.icons = icons;

ReactDOM.render(
  <Provider store={store}>
    <I18nextProvider i18n={i18next}>
      <App />
    </I18nextProvider>
  </Provider>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
