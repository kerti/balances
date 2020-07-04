import React from "react";
import { useTranslation } from "react-i18next";

const Properties = () => {
  const { t } = useTranslation("main");
  return (
    <>
      <div className="card">
        <div className="card-header">{t("assets.properties.properties")}</div>
        <div className="card-body">
          <p>{t("assets.properties.description")}</p>
        </div>
      </div>
    </>
  );
};

export default Properties;
