import React from "react";
import { TheContent, TheSidebar, TheFooter, TheHeader } from "./index";
import { useSelector } from "react-redux";
import { Redirect } from "react-router-dom";

const TheLayout = () => {
  const token = useSelector((state) => state.auth.token);
  const tokenExpiration = useSelector((state) => state.auth.expiration);
  const expirationDate = new Date(tokenExpiration);
  const expired = new Date() > expirationDate;
  return token && !expired ? (
    <div className="c-app c-default-layout">
      <TheSidebar />
      <div className="c-wrapper">
        <TheHeader />
        <div className="c-body">
          <TheContent />
        </div>
        <TheFooter />
      </div>
    </div>
  ) : (
    <Redirect to="/login" />
  );
};

export default TheLayout;
