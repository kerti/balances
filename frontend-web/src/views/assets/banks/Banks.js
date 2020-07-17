import React, { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useHistory } from "react-router-dom";
import {
  CRow,
  CCol,
  CCard,
  CCardHeader,
  CCardBody,
  CDataTable,
  CBadge,
  CPagination,
} from "@coreui/react";
import { useDispatch, useSelector } from "react-redux";
import { requestBankList } from "../../../data/actions/assets/banks";

const getBadge = (item) => {
  switch (item.status) {
    case "active":
      return "success";
    case "inactive":
      return "secondary";
    default:
      return "primary";
  }
};

const Banks = () => {
  const { t } = useTranslation(["assets", "formats"]);
  const { t: f } = useTranslation("formats");
  const dispatch = useDispatch();
  const [currentPage, setCurrentPage] = useState(1);
  const history = useHistory();
  const fields = [
    { key: "accountName", label: t("banks.accountName") },
    { key: "bankName", label: t("banks.name") },
    { key: "accountHolderName", label: t("banks.accountHolder") },
    { key: "accountNumber", label: t("banks.accountNumber") },
    {
      key: "lastBalance",
      label: t("banks.lastBalance"),
      _classes: ["text-right"],
    },
    { key: "status", label: t("banks.status") },
  ];

  const bankData = useSelector((state) => state.pages.assets.banks);

  useEffect(() => {
    if (currentPage > 0) {
      dispatch(
        requestBankList(
          "",
          currentPage,
          parseInt(process.env.REACT_APP_DEFAULT_PAGE_SIZE)
        )
      );
    }
  }, [dispatch, currentPage]);

  const spinner = (
    <div className="pt-3 text-center">
      <div className="spinner-border text-primary" role="status">
        <span className="sr-only">Loading...</span>
      </div>
    </div>
  );

  const dataTable = (
    <>
      <CDataTable
        items={bankData.items}
        fields={fields}
        hover
        striped
        bordered
        size="sm"
        itemsPerPage={parseInt(process.env.REACT_APP_DEFAULT_PAGE_SIZE)}
        clickableRows
        onRowClick={(item) => {
          history.push(`/assets/banks/${item.id}`);
        }}
        scopedSlots={{
          lastBalance: (item) => (
            <td className="text-right">
              {f("number.decimal.2fractions", {
                value: item.lastBalance,
              })}
              <br />
              <small>
                {f("date.long", {
                  value: item.lastBalanceDate,
                })}
              </small>
            </td>
          ),
          status: (item) => (
            <td>
              <CBadge color={getBadge(item)}>
                {t("banks.states." + item.status)}
              </CBadge>
            </td>
          ),
        }}
      />
      <CPagination
        size="sm"
        activePage={currentPage}
        pages={bankData.pageInfo ? bankData.pageInfo.totalPages : 0}
        onActivePageChange={setCurrentPage}
      />
    </>
  );

  return (
    <CRow>
      <CCol>
        <CCard>
          <CCardHeader>{t("banks.listOfBankAccounts")}</CCardHeader>
          <CCardBody>{bankData.loading ? spinner : dataTable}</CCardBody>
        </CCard>
      </CCol>
    </CRow>
  );
};

export default Banks;
