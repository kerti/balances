import React from "react";
import { useTranslation } from "react-i18next";

const Institutional = () => {
  const { t } = useTranslation("liabilities");
  return (
    <>
      <div className="card">
        <div className="card-header">{t("institutional.institutional")}</div>
        <div className="card-body">
          <p>{t("institutional.description")}</p>
        </div>
      </div>
    </>
  );
};

export default Institutional;
