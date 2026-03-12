let currentApiUrl = '';
let currentOffset = 0;
const LIMIT = 10;
let totalCount = 0;

/**
 * 1. 司令塔: パス解析とタブの初期化
 */
const dispatchFollowTask = async () => {
  const path = window.location.pathname;
  const pathParts = path.split('/');
  const userId = pathParts[2];

  const tabs = document.querySelectorAll('header div');
  const followingTab = tabs[0]; // 1つ目の要素：フォロー中
  const followerTab = tabs[1]; // 2つ目の要素：フォロワー

  if (path.includes('followings')) {
    currentApiUrl = `/api/users/${userId}/followings`;
    updateTabUI(followingTab, followerTab);
  } else if (path.includes('followers')) {
    currentApiUrl = `/api/users/${userId}/followers`;
    updateTabUI(followerTab, followingTab);
  }

  setupPagination();
  loadUsers();
};

/**
 * 2. タブの見た目の切り替え (アクティブな方に青い下線を出す)
 */
const updateTabUI = (activeTab, inactiveTab) => {
  // アクティブ側のスタイル設定
  activeTab.classList.remove('text-gray-500');
  activeTab.classList.add('font-bold', 'text-black');

  // 青い下線インジケーター（X/Twitter風）
  const indicator = document.createElement('div');
  indicator.className = 'absolute bottom-0 w-16 h-1 bg-[#1d9bf0] rounded-full';
  activeTab.appendChild(indicator);

  // 非アクティブ側のスタイル設定
  inactiveTab.classList.add('text-gray-500');
  inactiveTab.classList.remove('font-bold', 'text-black');
};

/**
 * 3. ユーザー一覧の取得
 */
const loadUsers = async (offset = 0) => {
  currentOffset = offset;
  try {
    const response = await fetch(
      `${currentApiUrl}?limit=${LIMIT}&offset=${currentOffset}`,
    );
    if (!response.ok) throw new Error('データ取得失敗');

    const data = await response.json();
    totalCount = data.count || 0;
    const users = data.follow_list || [];

    const container = document.getElementById('user-list');
    container.innerHTML = '';

    users.forEach((user) => {
      container.appendChild(createUserCard(user));
    });

    updatePaginationUI();
  } catch (error) {
    console.error('Error:', error);
  }
};

/**
 * 4. ユーザーカードの作成 (ツイート一覧と同じ形式)
 */
const createUserCard = (user) => {
  const card = document.createElement('div');
  // ツイート一覧と同じパディングとボーダー、ホバー効果を付与
  card.className =
    'p-4 border-b border-gray-100 hover:bg-gray-50/50 transition cursor-pointer flex gap-3';

  card.innerHTML = `
    <img src="${user.profile_image}" class="w-10 h-10 rounded-full bg-gray-200 flex-shrink-0 object-cover" alt="Avatar">
    
    <div class="flex-1">
      <div class="flex flex-col">
        <span class="font-bold text-[15px] hover:underline">${user.user_name}</span>
        <span class="text-gray-500 text-[14px]">@${user.user_name}</span>
      </div>
      <p class="text-[15px] leading-5 mt-1 text-gray-800">${user.self_introduction || ''}</p>
    </div>
  `;

  // カード全体をクリックしたらその人の詳細ページへ飛ぶようにする場合
  card.addEventListener('click', () => {
    window.location.href = `/user-detail/${user.id}`;
  });

  return card;
};

/**
 * 5. ページネーション制御
 */
const setupPagination = () => {
  document.getElementById('prev-btn').addEventListener('click', () => {
    if (currentOffset >= LIMIT) loadUsers(currentOffset - LIMIT);
  });
  document.getElementById('next-btn').addEventListener('click', () => {
    if (currentOffset + LIMIT < totalCount) loadUsers(currentOffset + LIMIT);
  });
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

// 実行開始
dispatchFollowTask();
