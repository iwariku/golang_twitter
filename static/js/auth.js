const signUp = () => {
  document.getElementById('signupForm')?.addEventListener('submit', (e) => {
    e.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    console.log('送信したデータ:', { email, password });

    fetch('/signup', {
      method: 'post',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });
  });
};

const login = () => {
  document.getElementById('loginForm')?.addEventListener('submit', (e) => {
    e.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    console.log('送信したデータ:', { email, password });

    fetch('/login', {
      method: 'post',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });
  });
};

const dispatchPathTask = () => {
  const path = window.location.pathname;
  if (path.includes('signup')) {
    signUp();
  } else if (path.includes('login')) {
    login();
  }
};

dispatchPathTask();
