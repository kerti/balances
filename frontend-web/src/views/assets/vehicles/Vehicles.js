import React from "react";
import { useTranslation } from "react-i18next";

const Vehicles = () => {
  const { t, i18n } = useTranslation("main");
  return (
    <>
      <div className="card">
        <div className="card-header">{t("assets.vehicles.vehicles")}</div>
        <div className="card-body">
          <p>{t("assets.vehicles.description")}</p>
        </div>
      </div>
    </>
  );
};

export default Vehicles;
