import { Link } from "react-router-dom";
import "./header.css";

function Header() {
  return (
    <header className="back-header-container">
      <Link to="/" className="back-header-link" aria-label="Back to Home">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth={2}
          stroke="currentColor"
          className="back-header-icon"
          aria-hidden="true"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M15 19l-7-7 7-7"
          />
        </svg>
        <span>Back Home</span>
      </Link>
    </header>
  );
}

export default Header