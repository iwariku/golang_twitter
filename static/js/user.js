const getUser = async () => {
  // 現在のパス (/user-detail/1) から最後の部分を取り出す
  const pathParts = window.location.pathname.split('/');
  const userId = pathParts[pathParts.length - 1]; // 一番後ろの "1" を取得

  if (!userId || isNaN(userId)) {
    console.error('ユーザーIDがURLに含まれていません');
    return;
  }

  const response = await fetch(`/api/users/${userId}`);
  const data = await response.json();

  console.log('ユーザー情報取得完了');
  container = document.getElementById('user-detail-container');
  container.innerHTML = `
  <div class="p-10">
    <img src="${data.profile_image}" class="w-24 h-24 rounded-full border-4 border-white -mt-12 ml-4 object-cover shadow-sm bg-gray-200" alt="プロフィール画像">
    <div class="mt-3 px-4">
      <div class="text-xl font-bold">${data.user_name}</div>
      <div class="text-gray-500 text-sm">@${data.user_name}</div>
      <div class="mt-3 text-gray-800">${data.self_introduction}</div>
      <div class="mt-3 text-gray-500 text-sm">
        <span>誕生日: ${data.date_of_birth}</span>
      </div>
    </div>
  </div>
`;
};
