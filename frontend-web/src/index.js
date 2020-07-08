import "react-app-polyfill/ie11"; // For IE 11 support
import "react-app-polyfill/stable";
import "./polyfill";
import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import * as serviceWorker from "./serviceWorker";

import { icons } from "./assets/icons";

import { Provider } from "react-redux";
import initStore from "./data/store";

import { CookiesProvider } from "react-cookie";

import { I18nextProvider } from "react-i18next";
import i18next from "i18next";

import app_en from "./translations/en/app.json";
import app_id from "./translations/id/app.json";

import assets_en from "./translations/en/assets.json";
import assets_id from "./translations/id/assets.json";

import liabilities_en from "./translations/en/liabilities.json";
import liabilities_id from "./translations/id/liabilities.json";

import investments_en from "./translations/en/investments.json";
import investments_id from "./translations/id/investments.json";

import navigation_en from "./translations/en/navigation.json";
import navigation_id from "./translations/id/navigation.json";

i18next.init({
  interpolation: { escapeValue: false },
  lng: process.env.REACT_APP_DEFAULT_LANG,
  resources: {
    en: {
      app: app_en,
      assets: assets_en,
      liabilities: liabilities_en,
      investments: investments_en,
      navigation: navigation_en,
    },
    id: {
      app: app_id,
      assets: assets_id,
      liabilities: liabilities_id,
      investments: investments_id,
      navigation: navigation_id,
    },
  },
});

React.icons = icons;

const store = initStore();

ReactDOM.render(
  <CookiesProvider>
    <Provider store={store}>
      <I18nextProvider i18n={i18next}>
        <App />
      </I18nextProvider>
    </Provider>
  </CookiesProvider>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
