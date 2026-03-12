// ユーザーAPIを叩く処理
const getUser = async () => {
  const pathParts = window.location.pathname.split('/');
  const userId = pathParts[pathParts.length - 1];

  if (!userId || isNaN(userId)) {
    console.error('ユーザーIDがURLに含まれていません');
    return;
  }

  try {
    const response = await fetch(`/api/users/${userId}`);
    if (!response.ok) throw new Error('ユーザー情報の取得に失敗しました');
    const data = await response.json();

    setupUserInfo(data);

    // フォローボタンとラベルのセットアップ
    const followBtn = document.getElementById('js-follow-btn');
    const followingCountLable = document.getElementById('following-count');
    const followerCountLabel = document.getElementById('js-follower-count');

    setupFollowButton(followBtn, userId, followerCountLabel);
  } catch (error) {
    console.error('getUser Error:', error);
  }
};

// ユーザー情報の形成
const setupUserInfo = (data) => {
  const container = document.getElementById('user-detail-container');

  // Tailwind CSSを使用して、画像のUIを再現
  container.innerHTML = `
    <div class="relative">
      <div class="h-32 bg-gray-200"></div>

      <div class="px-4 pb-4">
        <div class="relative flex justify-between items-start">
          <img src="${data.profile_image}" 
              class="w-24 h-24 rounded-full border-4 border-white -mt-12 object-cover bg-gray-100 shadow-sm" 
              alt="プロフィール画像">
          
          <div class="mt-3">
            <button id="js-follow-btn" data-is-followed="${data.is_followed}"></button>
          </div>
        </div>

        <div class="mt-3">
          <div class="text-xl font-extrabold leading-tight">${data.user_name}</div>
          <div class="text-gray-500 text-[15px]">@${data.user_name}</div>
          
          <div class="mt-3 text-[15px] leading-normal text-gray-900">
            ${data.self_introduction}
          </div>

          <div class="mt-3 text-gray-500 text-sm flex items-center gap-1">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
            </svg>
            <span>誕生日: ${data.date_of_birth}</span>
          </div>

          <div class="mt-3 flex gap-4 text-sm">
            <a href="/" class="hover:underline">
              <span id="js-following-count" class="font-bold text-black">${data.following_count}</span>
              <span class="text-gray-500 font-normal">フォロー中</span>
            </a>
            <a href="/" class="hover:underline">
              <span id="js-follower-count" class="font-bold text-black">${data.follower_count}</span>
              <span class="text-gray-500 font-normal">フォロワー</span>
            </a>
          </div>
        </div>
      </div>
    </div>
  `;
};

// ボタンを押して状態が変わるまでの動き
const setupFollowButton = (btn, targetUserId, countLabel) => {
  if (!btn) return;

  applyFollowButtonStyle(btn);

  btn.addEventListener('click', () => {
    actionFollow(btn, targetUserId, countLabel);
  });
};

// フォローボタンを押した時の処理
const actionFollow = async (btn, targetUserId, countLabel) => {
  const isFollowed = btn.dataset.isFollowed === 'true';
  const method = isFollowed ? 'DELETE' : 'POST';

  try {
    const response = await fetch(`/api/users/${targetUserId}/follow`, {
      method: method,
    });

    if (!response.ok) throw new Error('サーバーエラー');

    const data = await response.json();
    btn.dataset.isFollowed = data.is_followed;

    applyFollowButtonStyle(btn);
  } catch (error) {
    console.error('フォロー操作に失敗しました');
  }
};

// ボタンのUI
const applyFollowButtonStyle = (btn) => {
  const isFollowed = btn.dataset.isFollowed === 'true';

  if (isFollowed) {
    btn.className =
      'px-5 py-1.5 rounded-full font-bold text-sm transition bg-white text-black border border-gray-300 hover:bg-red-50 hover:text-red-600 hover:border-red-200';
    btn.textContent = 'フォロー解除';
  } else {
    btn.className =
      'px-5 py-1.5 rounded-full font-bold text-sm transition bg-black text-white hover:bg-gray-800';
    btn.textContent = 'フォロー';
  }
};
