const createGroup = async (name) => {
  const response = await fetch('/api/dm/group', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name }),
  });
  if (!response.ok) throw new Error(`作成失敗: ${response.status}`);
  return await response.json();
};

const addMember = async (userId, groupId) => {
  const response = await fetch(`/api/dm/add-member`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ user_id: userId, group_id: groupId }),
  });
  if (!response.ok) throw new Error(`追加失敗: ${response.status}`);
  return await response.json();
};

const sendMessage = async (groupId, message) => {
  const response = await fetch(`/api/dm/groups/${groupId}/message`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ message }),
  });
  if (!response.ok) throw new Error(`送信失敗: ${response.status}`);
  return await response.json();
};

const loadMessages = async (groupId) => {
  const container = document.getElementById('messageList');
  if (!container) return;

  const response = await fetch(`/api/dm/groups/${groupId}/messages`);
  const data = await response.json();
  const messages = data.messages || [];

  container.innerHTML = messages.length
    ? ''
    : '<div class="p-4 text-center text-gray-500">メッセージがありません</div>';

  messages.forEach((msg) => {
    const card = document.createElement('div');
    card.className =
      'p-4 border-b border-gray-100 hover:bg-gray-50 transition flex flex-col gap-1';
    card.innerHTML = `
      <div class="text-[12px] text-gray-500 font-bold">User ID: ${msg.user_id}</div>
      <div class="text-[15px] text-gray-800">${msg.message}</div>
    `;
    container.appendChild(card);
  });
};

const loadGroups = async () => {
  const container = document.getElementById('groupList');
  if (!container) return;

  const response = await fetch('/api/dm/groups');
  const data = await response.json();
  const groups = data.groups || [];

  container.innerHTML = groups.length
    ? ''
    : '<div class="p-10 text-center text-gray-500">参加中のグループはありません</div>';

  groups.forEach((group) => {
    const card = document.createElement('div');
    card.className =
      'p-4 border-b border-gray-100 hover:bg-gray-50 transition cursor-pointer flex items-center gap-4';
    card.innerHTML = `
      <div class="w-12 h-12 bg-gray-200 rounded-full flex-shrink-0 flex items-center justify-center">👥</div>
      <div class="flex-1">
        <div class="flex justify-between items-center"><span class="font-bold">${group.name}</span></div>
        <p class="text-sm text-gray-500 mt-1">メッセージを見る</p>
      </div>
    `;
    card.onclick = () =>
      (window.location.href = `/dm/groups/${group.dm_group_id}/messages`);
    container.appendChild(card);
  });
};

const setupCreateGroupForm = () => {
  document
    .getElementById('createGroupForm')
    ?.addEventListener('submit', async (e) => {
      e.preventDefault();
      const name = document.getElementById('name').value;
      try {
        await createGroup(name);
        alert('作成成功！');
        window.location.href = '/dm/add-member';
      } catch (err) {
        alert(err.message);
      }
    });
};

const setupAddMemberForm = () => {
  document
    .getElementById('addMemberForm')
    ?.addEventListener('submit', async (e) => {
      e.preventDefault();
      const userId = parseInt(document.getElementById('user_id').value);
      const groupId = parseInt(document.getElementById('group_id').value);
      try {
        await addMember(userId, groupId);
        alert('メンバーを追加しました');
        window.location.href = '/dm/groups';
      } catch (err) {
        alert(err.message);
      }
    });
};

const setupSendMessageForm = (groupId) => {
  document
    .getElementById('sendMessageForm')
    ?.addEventListener('submit', async (e) => {
      e.preventDefault();
      const message = document.getElementById('messageContent').value;
      try {
        await sendMessage(groupId, message);
        alert('送信しました！');
        window.location.href = `/dm/groups/${groupId}/messages`;
      } catch (err) {
        alert(err.message);
      }
    });
};

const dispatchDmTask = async () => {
  const path = window.location.pathname;
  const pathParts = path.split('/');

  // /dm/groups/:id/messages
  if (path.includes('/messages')) {
    const groupId = pathParts[3];
    loadMessages(groupId);
  }
  // /dm/groups/:id/message (送信画面)
  else if (path.endsWith('/message')) {
    const groupId = pathParts[3];
    setupSendMessageForm(groupId);
  }
  // /dm/groups (一覧)
  else if (path === '/dm/groups') {
    loadGroups();
  }
  // /dm/group (作成)
  else if (path === '/dm/group') {
    setupCreateGroupForm();
  }
  // /dm/add-member (追加)
  else if (path === '/dm/add-member') {
    setupAddMemberForm();
  }
};

document.addEventListener('DOMContentLoaded', dispatchDmTask);
