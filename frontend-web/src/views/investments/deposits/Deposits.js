import React from "react";
import { useTranslation } from "react-i18next";

const Deposits = () => {
  const { t } = useTranslation("investments");
  return (
    <>
      <div className="card">
        <div className="card-header">{t("deposits.deposits")}</div>
        <div className="card-body">
          <p>{t("deposits.description")}</p>
        </div>
      </div>
    </>
  );
};

export default Deposits;
