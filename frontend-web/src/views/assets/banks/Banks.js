import React from "react";
import { useTranslation } from "react-i18next";

const Banks = () => {
  const { t, i18n } = useTranslation("main");
  return (
    <>
      <div className="card">
        <div className="card-header">{t("assets.banks.bankAccounts")}</div>
        <div className="card-body">
          <p>{t("assets.banks.description")}</p>
        </div>
      </div>
    </>
  );
};

export default Banks;
