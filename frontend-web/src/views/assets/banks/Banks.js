import React, { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useDispatch } from "react-redux";
import { loadBankAccountPage } from "../../../data/actions/assets/bankAccounts";

const Banks = () => {
  const { t } = useTranslation("assets");
  const dispatch = useDispatch();
  console.log("banks page rendered");

  useEffect(() => {
    dispatch(loadBankAccountPage("", 1));
  }, [dispatch]);

  const nextPage = () => {
    dispatch(loadBankAccountPage("", 2));
  };

  return (
    <>
      <div className="card">
        <div className="card-header">{t("banks.bankAccounts")}</div>
        <div className="card-body" onClick={nextPage}>
          <p>{t("banks.description")}</p>
        </div>
      </div>
    </>
  );
};

export default Banks;
