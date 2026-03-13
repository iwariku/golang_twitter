let currentApiUrl = '';
let currentOffset = 0;
const LIMIT = 3;
let totalCount = 0;

// パスによって叩くAPIを変える
const dispatchFollowTask = async () => {
  const path = window.location.pathname;
  const PathParts = path.split('/');
  const userId = PathParts[2];

  const titleElem = document.getElementById('page-title');

  if (path.includes('followings')) {
    titleElem.textContent = 'フォロー中';
    currentApiUrl = `/api/users/${userId}/followings`;
  } else if (path.includes('followers')) {
    titleElem.textContent = 'フォロワー';
    currentApiUrl = `/api/users/${userId}/followers`;
  }

  setupPagination();
  loadUsers();
};

// ユーザー一覧の取得
const loadUsers = async (offset = 0) => {
  currentOffset = offset;
  try {
    const response = await fetch(
      `${currentApiUrl}?limit=${LIMIT}&offset=${currentOffset}`,
    );
    if (!response.ok) throw new Error('データ取得失敗');

    const data = await response.json();
    totalCount = data.count;
    const users = data.follow_list;

    const container = document.getElementById('user-list');
    container.innerHTML = ``;

    if (users.length === 0) {
      container.innerHTML`<div class="text-center text-gray-500 font-bold">表示するユーザーがいません</div>`;
    }

    users.forEach((user) => {
      container.appendChild(createUserCard(user));
    });
    updatePaginationUI();
  } catch (error) {
    console.error('Error', error);
  }
};

// ユーザーカードの作成
const createUserCard = (user) => {
  const card = document.createElement('div');
  card.className =
    'p-4 border-b border-gray-100 hover:bg-gray-50/50 transition cursor-pointer flex gap-3';
  card.innerHTML = `
    <img src="${user.profile_image || '/static/images/default-avatar.png'}" class="w-10 h-10 rounded-full bg-gray-200 flex-shrink-0 object-cover" alt="Avatar">
    <div class="flex-1">
      <div class="flex flex-col">
        <span class="font-bold text-[15px] hover:underline">${user.user_name}</span>
        <span class="text-gray-500 text-[14px]">@${user.user_name}</span>
      </div>
      <p class="text-[15px] leading-5 mt-1 text-gray-800">${user.self_introduction || ''}</p>
    </div>
  `;

  // カードクリックで詳細画面へ
  card.onclick = () => {
    window.location.href = `/user-detail/${user.id}`;
  };

  return card;
};

// 4. ページネーション制御(ツイート一覧と一緒)
const setupPagination = () => {
  document.getElementById('prev-btn').onclick = () => {
    if (currentOffset >= LIMIT) loadUsers(currentOffset - LIMIT);
  };
  document.getElementById('next-btn').onclick = () => {
    if (currentOffset + LIMIT < totalCount) loadUsers(currentOffset + LIMIT);
  };
};

const updatePaginationUI = () => {
  const currentPage = Math.floor(currentOffset / LIMIT) + 1;
  const maxPage = Math.ceil(totalCount / LIMIT) || 1;
  document.getElementById('page-info').textContent =
    `${currentPage} / ${maxPage} ページ (全 ${totalCount} 件)`;
  document.getElementById('prev-btn').disabled = currentOffset === 0;
  document.getElementById('next-btn').disabled =
    currentOffset + LIMIT >= totalCount;
};

// 実行
dispatchFollowTask();
