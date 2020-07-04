import React from "react";
import { useTranslation } from "react-i18next";

const Personal = () => {
  const { t } = useTranslation("liabilities");
  return (
    <>
      <div className="card">
        <div className="card-header">{t("personal.personal")}</div>
        <div className="card-body">
          <p>{t("personal.description")}</p>
        </div>
      </div>
    </>
  );
};

export default Personal;
