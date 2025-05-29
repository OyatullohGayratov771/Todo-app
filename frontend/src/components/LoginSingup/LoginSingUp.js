import React from "react";
import "./LoginSingUp.css";

import user_icon from "../Assets/user.png";
import mail_icon from "../Assets/mail.png";
import lock_icon from "../Assets/locked-computer.png";

const LoginSingUp = () => {
    return (
        <div className="container">
            <div className="header">
                <div className="text">Sing Up</div>
                <div className="underline"></div>
            </div>
                <div className="input">
                    <img src={user_icon} alt="" />
                    <input type="text" />
                </div>
            <div className="inputs">
                <div className="input">
                    <img src={mail_icon} alt="" />
                    <input type="email" />
                </div>
                <div className="input">
                    <img src={lock_icon} alt="" />
                    <input type="password" />
                </div>
            </div>
        </div>
    );
};

export default LoginSingUp;