import React from "react";
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
} from "@coreui/react";

// TODO: remove this and use actual data
const bankAccountData = [
  {
    id: "0c2fa766-6726-4513-aa4d-fafe3010d114",
    accountName: "Current",
    bankName: "First Bank of Texas",
    accountHolder: "Leonardo da Vinci",
    accountNumber: "634.545.6.4526.23600",
    lastBalance: 4234920.55,
    lastBalanceDate: 1592223589000,
    status: "active",
  },
  {
    id: "75e34690-5505-4587-8d4b-caeca7fa24b2",
    accountName: "Savings",
    bankName: "Standard Chartered London",
    accountHolder: "Michelangelo di Lodovico",
    accountNumber: "24523-642.34353-001",
    lastBalance: 365324525.33,
    lastBalanceDate: 1594023589000,
    status: "active",
  },
  {
    id: "c468f41e-0f45-4f65-a10e-0001e85c8347",
    accountName: "Emergency Fund",
    bankName: "HSBC Bangkok",
    accountHolder: "Donato di Niccolò",
    accountNumber: "69392424642343",
    lastBalance: 90532740.23,
    lastBalanceDate: 1534023589000,
    status: "inactive",
  },
  {
    id: "f68fc8cc-ffd5-42a9-9f55-732f68c4ee1b",
    accountName: "Investment Pool",
    bankName: "Commonwealth Bank of Australia",
    accountHolder: "Raffaello Sanzio da Urbino",
    accountNumber: "534234034",
    lastBalance: 0,
    lastBalanceDate: 1594423589000,
    status: "inactive",
  },
];

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
  const history = useHistory();
  const fields = [
    { key: "accountName", label: t("banks.accountName") },
    { key: "bankName", label: t("banks.name") },
    { key: "accountHolder", label: t("banks.accountHolder") },
    { key: "accountNumber", label: t("banks.accountNumber") },
    {
      key: "lastBalance",
      label: t("banks.lastBalance"),
      _classes: ["text-right"],
    },
    { key: "status", label: t("banks.status") },
  ];

  return (
    <>
      <CRow>
        <CCol>
          <CCard>
            <CCardHeader>{t("banks.listOfBankAccounts")}</CCardHeader>
            <CCardBody>
              <CDataTable
                items={bankAccountData}
                fields={fields}
                hover
                striped
                bordered
                size="sm"
                itemsPerPage={10}
                pagination
                clickableRows
                onRowClick={(item) => {
                  console.log("row clicked");
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
            </CCardBody>
          </CCard>
        </CCol>
      </CRow>
    </>
  );
};

export default Banks;
