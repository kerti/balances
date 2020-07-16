import React, { useEffect } from "react";
import { useTranslation } from "react-i18next";
import {
  CRow,
  CCol,
  CCard,
  CCardHeader,
  CCardBody,
  CForm,
  CFormGroup,
  CLabel,
  CInput,
  CCardFooter,
  CButton,
  CDataTable,
} from "@coreui/react";
import CIcon from "@coreui/icons-react";
import { CChartLine } from "@coreui/react-chartjs";
import { useDispatch, useSelector } from "react-redux";
import { requestBankDetailByID } from "../../../data/actions/assets/banks";
import { useParams } from "react-router-dom";

// TODO: remove this and use actual data
const balanceHistoryData = [
  {
    id: "20a4b907-4eb5-4591-93e6-e439cb75e224",
    bankAccountId: "0c2fa766-6726-4513-aa4d-fafe3010d114",
    date: 1582028589000,
    balance: 459493242.23,
    setBy: "Rafaello Sanzio da Urbino",
  },
  {
    id: "20a4b907-4eb5-4591-93e6-e439cb75e225",
    bankAccountId: "0c2fa766-6726-4513-aa4d-fafe3010d114",
    date: 1583028589000,
    balance: 459493242.23,
    setBy: "Rafaello Sanzio da Urbino",
  },
  {
    id: "20a4b907-4eb5-4591-93e6-e439cb75e226",
    bankAccountId: "0c2fa766-6726-4513-aa4d-fafe3010d114",
    date: 1584028589000,
    balance: 459493242.23,
    setBy: "Leonardo da Vinci",
  },
  {
    id: "20a4b907-4eb5-4591-93e6-e439cb75e227",
    bankAccountId: "0c2fa766-6726-4513-aa4d-fafe3010d114",
    date: 1585028589000,
    balance: 459493242.23,
    setBy: "Rafaello Sanzio da Urbino",
  },
  {
    id: "20a4b907-4eb5-4591-93e6-e439cb75e228",
    bankAccountId: "0c2fa766-6726-4513-aa4d-fafe3010d114",
    date: 1590028589000,
    balance: 454493242.23,
    setBy: "Donato di Nicolò",
  },
  {
    id: "20a4b907-4eb5-4591-93e6-e439cb75e229",
    bankAccountId: "0c2fa766-6726-4513-aa4d-fafe3010d114",
    date: 1592028589000,
    balance: 455493242.23,
    setBy: "Michelangelo di Lodovico",
  },
  {
    id: "20a4b907-4eb5-4591-93e6-e439cb75e230",
    bankAccountId: "0c2fa766-6726-4513-aa4d-fafe3010d114",
    date: 1593028589000,
    balance: 456493242.23,
    setBy: "Michelangelo di Lodovico",
  },
  {
    id: "20a4b907-4eb5-4591-93e6-e439cb75e231",
    bankAccountId: "0c2fa766-6726-4513-aa4d-fafe3010d114",
    date: 1594028589000,
    balance: 457493242.23,
    setBy: "Donato di Nicolò",
  },
  {
    id: "20a4b907-4eb5-4591-93e6-e439cb75e232",
    bankAccountId: "0c2fa766-6726-4513-aa4d-fafe3010d114",
    date: 1595028589000,
    balance: 451493242.23,
    setBy: "Leonardo da Vinci",
  },
  {
    id: "20a4b907-4eb5-4591-93e6-e439cb75e233",
    bankAccountId: "0c2fa766-6726-4513-aa4d-fafe3010d114",
    date: 1596028589000,
    balance: 458493242.23,
    setBy: "Donato di Nicolò",
  },
];

const getDatasetFromBalanceHistory = (balanceHistoryData) => {
  return balanceHistoryData.map((hist) => ({
    x: new Date(hist.date),
    y: hist.balance,
  }));
};

const getLabelsFromBalanceHistory = (balanceHistoryData) => {
  return balanceHistoryData.map((hist) => new Date(hist.date));
};

const chartOptions = {
  scales: {
    xAxes: [
      {
        type: "time",
        distribution: "linear",
      },
    ],
  },
  tooltips: {
    enabled: true,
  },
};

const Properties = () => {
  const { t } = useTranslation("assets");
  const { t: f } = useTranslation("formats");
  const dispatch = useDispatch();
  const balanceHistoryFields = [
    { key: "date", label: t("common.date") },
    { key: "setBy", label: t("banks.balanceSetBy") },
    { key: "balance", label: t("banks.balance"), _classes: ["text-right"] },
  ];

  const { id } = useParams();
  const bankData = useSelector((state) => state.pages.assets.banks);
  console.log("ASDASDASDASD");
  console.log(bankData);

  useEffect(() => {
    dispatch(requestBankDetailByID(id, true, 10));
  }, [dispatch, id]);
  return (
    <>
      <CRow>
        <CCol xs="12" md="12">
          <CCard>
            <CCardHeader>{t("banks.details")}</CCardHeader>
            <CCardBody>
              <CForm action="" method="post" className="form-horizontal">
                <CFormGroup row>
                  <CCol md="3">
                    <CLabel htmlFor="bank-name">
                      {t("banks.accountName")}
                    </CLabel>
                  </CCol>
                  <CCol xs="12" md="9">
                    <CInput
                      type="text"
                      id="account-name"
                      name="account-name"
                      placeholder={t("banks.accountNamePlaceholder")}
                    />
                  </CCol>
                </CFormGroup>
                <CFormGroup row>
                  <CCol md="3">
                    <CLabel htmlFor="bank-name">{t("banks.name")}</CLabel>
                  </CCol>
                  <CCol xs="12" md="9">
                    <CInput
                      type="text"
                      id="bank-name"
                      name="bank-name"
                      placeholder={t("banks.namePlaceholder")}
                    />
                  </CCol>
                </CFormGroup>
                <CFormGroup row>
                  <CCol md="3">
                    <CLabel htmlFor="account-holder">
                      {t("banks.accountHolder")}
                    </CLabel>
                  </CCol>
                  <CCol xs="12" md="9">
                    <CInput
                      type="text"
                      id="account-holder"
                      name="account-holder"
                      placeholder={t("banks.accountHolderPlaceholder")}
                    />
                  </CCol>
                </CFormGroup>
                <CFormGroup row>
                  <CCol md="3">
                    <CLabel htmlFor="account-number">
                      {t("banks.accountNumber")}
                    </CLabel>
                  </CCol>
                  <CCol xs="12" md="9">
                    <CInput
                      type="text"
                      id="account-number"
                      name="account-number"
                      placeholder={t("banks.accountNumberPlaceholder")}
                    />
                  </CCol>
                </CFormGroup>
              </CForm>
            </CCardBody>
            <CCardFooter>
              <CButton type="submit" size="sm" color="primary">
                <CIcon name="cil-scrubber" /> {t("common.submit")}
              </CButton>{" "}
              <CButton type="reset" size="sm" color="danger">
                <CIcon name="cil-ban" /> {t("common.reset")}
              </CButton>
            </CCardFooter>
          </CCard>
        </CCol>
      </CRow>
      <CRow>
        <CCol xs="12" md="6">
          <CCard>
            <CCardHeader>{t("banks.balanceHistory")}</CCardHeader>
            <CCardBody>
              <CDataTable
                items={balanceHistoryData}
                fields={balanceHistoryFields}
                hover
                striped
                bordered
                size="sm"
                itemsPerPage={10}
                pagination
                scopedSlots={{
                  date: (item) => (
                    <td>{f("date.long", { value: item.date })}</td>
                  ),
                  balance: (item) => (
                    <td className="text-right">
                      {f("currency", { value: item.balance })}
                    </td>
                  ),
                }}
              ></CDataTable>
            </CCardBody>
          </CCard>
        </CCol>
        <CCol xs="12" md="6">
          <CCard>
            <CCardHeader>{t("banks.balanceHistoryGraph")}</CCardHeader>
            <CCardBody>
              <CChartLine
                type="line"
                datasets={[
                  {
                    label: t("banks.balanceHistory"),
                    backgroundColor: "rgb(0,156,195,0.6)",
                    data: getDatasetFromBalanceHistory(balanceHistoryData),
                  },
                ]}
                options={chartOptions}
                labels={getLabelsFromBalanceHistory(balanceHistoryData)}
              />
            </CCardBody>
          </CCard>
        </CCol>
      </CRow>
    </>
  );
};

export default Properties;
