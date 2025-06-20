import { Link } from 'react-router-dom';

function Home() {
  return (
    <>
      {/* Background cubes */}
      <div className="background" aria-hidden="true">
        <div className="cube"></div>
        <div className="cube"></div>
        <div className="cube"></div>
        <div className="cube"></div>
        <div className="cube"></div>
        <div className="cube"></div>
      </div>

      {/* Content */}
      <main className="container">
        <div className="logo fadeInUp" style={{ animationDelay: "0.3s" }}>
          rent-computing
        </div>

        <h1 className="fadeInUp" style={{ animationDelay: "0.6s" }}>
          Rent Powerful Remote Virtual Machines Instantly
        </h1>

        <p className="fadeInUp" style={{ animationDelay: "0.9s" }}>
          Deploy, connect, and compute on demand from anywhere in the world. Your
          cloud workstation, ready in seconds.
        </p>

        <div className="btn-group fadeInUp" style={{ animationDelay: "1.2s" }}>
          <Link to="/login" className="btn">
            Log In
          </Link>
          <Link to="/register" className="btn secondary">
            Register
          </Link>
        </div>
      </main>
    </>
  );
}

export default Home