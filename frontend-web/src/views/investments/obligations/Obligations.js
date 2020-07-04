import React from "react";
import { useTranslation } from "react-i18next";

const Obligations = () => {
  const { t } = useTranslation("investments");
  return (
    <>
      <div className="card">
        <div className="card-header">{t("obligations.obligations")}</div>
        <div className="card-body">
          <p>{t("obligations.description")}</p>
        </div>
      </div>
    </>
  );
};

export default Obligations;
