import "./header.css";

function Header({ onSignOut, username }) {
  return (
    <header className="auth-header">
      <div className="auth-logo">rent-computing</div>
      <div className="auth-profile">
        <span>Welcome, {username}!</span>
        <button className="auth-signout" onClick={onSignOut}>
          Sign Out
        </button>
      </div>
    </header>
  );
}

export default Header