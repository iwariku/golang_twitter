const createGroup = async () => {
  const nameInput = document.getElementById('name');
  const groupName = nameInput.value;

  try {
    const response = await fetch('/api/dm/group', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: groupName }),
    });

    if (!response.ok) {
      throw new Error(`作成失敗: ${response.status}`);
    }

    const data = await response.json();
    console.log('作成成功:', data);
    alert('作成成功！');

    // 成功したら一覧へ
    window.location.href = '/dm/groups';
  } catch (error) {
    console.error('エラー発生:', error);
    alert(error.message);
  }
};

document.getElementById('createGroupForm')?.addEventListener('submit', (e) => {
  e.preventDefault();
  createGroup();
});
