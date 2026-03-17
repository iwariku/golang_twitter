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

const addMemberToGroup = async (userId, groupId) => {
  try {
    const response = await fetch(`/api/dm/add-member`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ user_id: userId, group_id: groupId }),
    });

    if (!response.ok) throw new Error(`追加失敗: ${response.status}`);
    console.log('メンバー追加成功');
  } catch (error) {
    console.error('メンバー追加エラー:', error);
  }
};

document
  .getElementById('addMemberForm')
  ?.addEventListener('submit', async (e) => {
    e.preventDefault();

    const userIdInput = document.getElementById('user_id');
    const userId = parseInt(userIdInput.value);

    const groupIdInput = document.getElementById('group_id');
    const groupId = parseInt(groupIdInput.value);

    if (!groupId || isNaN(userId)) {
      alert('グループIDまたはユーザーIDが正しくありません');
      return;
    }

    await addMemberToGroup(userId, groupId);
  });
