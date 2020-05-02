import React from "react";
import { Navbar, Nav, Button, Dropdown } from "react-bootstrap";
import { connect } from "react-redux";
import SignUpModal from "./SignUpModal";
import SignInModal from "./SignInModal";
import { toast } from "react-toastify";
import { LOGOUT } from "../../../redux/actionTypes";
import { useHistory, Link } from "react-router-dom";

function Header(props) {
  const history = useHistory();

  function handleLogout() {
    props.onLogout();
    history.push("/");
    toast.success("Logged out.");
  }

  function handleAdminPanelClick() {
    history.push("/admin");
  }

  return (
    <>
      <Navbar bg="light">
        <Navbar.Brand>
          <Link to="/">Symposium</Link>
        </Navbar.Brand>
        <Nav className="justify-content-end" style={{ width: "100%" }}>
          {props.username.length > 0 && props.token.length > 0 ? (
            <Dropdown drop="left">
              <Dropdown.Toggle variant="primary">
                <img
                  src={
                    process.env.API_URL + "/avatars/" + props.userId + ".jpg"
                  }
                  width="25"
                  height="25"
                  style={{ borderRadius: "50%" }}
                />
                &nbsp;
                {props.username}
              </Dropdown.Toggle>

              <Dropdown.Menu>
                {props.access == "99" && (
                  <Dropdown.Item onClick={handleAdminPanelClick}>
                    Admin panel
                  </Dropdown.Item>
                )}
                <Dropdown.Item onClick={handleLogout}>Logout</Dropdown.Item>
              </Dropdown.Menu>
            </Dropdown>
          ) : (
            <>
              <SignUpModal></SignUpModal>
              &nbsp;
              <SignInModal></SignInModal>
            </>
          )}
        </Nav>
      </Navbar>
    </>
  );
}

function mapStateToProps(state) {
  return {
    userId: state.user.id,
    username: state.user.username,
    token: state.user.token,
    access: state.user.access,
  };
}

const mapDispatchToProps = (dispatch) => ({
  onLogout: () => dispatch({ type: LOGOUT }),
});

export default connect(mapStateToProps, mapDispatchToProps)(Header);
