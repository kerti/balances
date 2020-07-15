import React from "react";
import { TheContent, TheSidebar, TheFooter, TheHeader } from "./index";
import { Redirect } from "react-router-dom";
import Cookies from "cookies-js";
import cookieNames from "../data/cookies";

const TheLayout = () => {
  const token = Cookies.get(cookieNames.auth.token);
  return token ? (
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
