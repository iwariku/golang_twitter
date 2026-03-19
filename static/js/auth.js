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

const unsubscribe = async () => {
  document
    .getElementById('unsubscribeForm')
    ?.addEventListener('click', async () => {
      const response = await fetch('/api/user/unsubscribe', {
        method: 'DELETE',
      });
      const data = await response.json();
      alert('退会に成功しました');
    });
};

const dispatchPathTask = () => {
  const path = window.location.pathname;
  if (path.includes('signup')) {
    signUp();
  } else if (path.includes('login')) {
    login();
  } else if (path.includes('unsubscribe')) {
    unsubscribe();
  }
};

dispatchPathTask();
