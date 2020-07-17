import i18next from "i18next";

import app_en from "./en/app.json";
import app_id from "./id/app.json";

import assets_en from "./en/assets.json";
import assets_id from "./id/assets.json";

import formats from "./formats.json";

import liabilities_en from "./en/liabilities.json";
import liabilities_id from "./id/liabilities.json";

import investments_en from "./en/investments.json";
import investments_id from "./id/investments.json";

import navigation_en from "./en/navigation.json";
import navigation_id from "./id/navigation.json";

export const i18nResources = {
  en: {
    app: app_en,
    assets: assets_en,
    formats: formats,
    liabilities: liabilities_en,
    investments: investments_en,
    navigation: navigation_en,
  },
  id: {
    app: app_id,
    assets: assets_id,
    formats: formats,
    liabilities: liabilities_id,
    investments: investments_id,
    navigation: navigation_id,
  },
};

const formatters = {
  en: {
    "date.long": new Intl.DateTimeFormat("en", {
      year: "numeric",
      month: "long",
      day: "numeric",
    }),
    "number.decimal.default": new Intl.NumberFormat("en", {
      style: "decimal",
    }),
    "number.decimal.2fractions": new Intl.NumberFormat("en", {
      style: "decimal",
      minimumFractionDigits: 2,
    }),
  },
  id: {
    "date.long": new Intl.DateTimeFormat("id", {
      year: "numeric",
      month: "long",
      day: "numeric",
    }),
    "number.decimal.default": new Intl.NumberFormat("id", {
      style: "decimal",
    }),
    "number.decimal.2fractions": new Intl.NumberFormat("id", {
      style: "decimal",
      minimumFractionDigits: 2,
    }),
  },
};

export const getFormat = (value, format, lang) => {
  return formatters[lang][format].format(value);
};

export const initTranslations = () => {
  i18next.init({
    interpolation: {
      escapeValue: false,
      format: getFormat,
    },
    lng: process.env.REACT_APP_DEFAULT_LANG,
    resources: i18nResources,
  });
};
