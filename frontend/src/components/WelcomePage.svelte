<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { GetText } from "../../wailsjs/go/main/App";
  import logoUrl from "../assets/images/logo_transparent.png"; // 新增這行來匯入圖片

  const dispatch = createEventDispatcher();

  // 歡迎頁面選項
  interface WelcomeOption {
    id: string;
    title: string;
    description: string;
    icon: string;
    action: () => void;
  }
  let options: WelcomeOption[] = [];
  let appTitle = "歡迎使用 INSYRA Insights"; // 修改標題
  let appSubtitle = "資料分析與視覺化工具";
  let recentFilesTitle = "最近使用的檔案";
  let noRecentFiles = "尚無最近使用的檔案";

  // i18n 輔助函數
  async function t(key: string): Promise<string> {
    try {
      return await GetText(key);
    } catch (err) {
      console.warn(`翻譯鍵值 "${key}" 不存在，返回預設值`);
      return key;
    }
  }
  // 初始化選項和文字
  async function initOptions() {
    // 載入多語言文字
    appTitle = (await t("welcome.title")) || "歡迎使用 INSYRA Insights"; // 修改標題後備文字
    appSubtitle = (await t("welcome.subtitle")) || "資料分析與視覺化工具";
    recentFilesTitle = (await t("welcome.recent_files")) || "最近使用的檔案";
    noRecentFiles =
      (await t("welcome.no_recent_files")) || "尚無最近使用的檔案";

    options = [
      {
        id: "open_csv",
        title: (await t("welcome.open_csv")) || "開啟 CSV 檔案",
        description:
          (await t("welcome.open_csv_desc")) || "從 CSV 檔案匯入資料",
        icon: "📊",
        action: () => dispatch("action", { type: "open_csv" }),
      },
      {
        id: "open_json",
        title: (await t("welcome.open_json")) || "開啟 JSON 檔案",
        description:
          (await t("welcome.open_json_desc")) || "從 JSON 檔案匯入資料",
        icon: "🔗",
        action: () => dispatch("action", { type: "open_json" }),
      },
      {
        id: "open_sqlite",
        title: (await t("welcome.open_sqlite")) || "開啟 SQLite 資料庫",
        description:
          (await t("welcome.open_sqlite_desc")) || "連接到 SQLite 資料庫",
        icon: "🗄️",
        action: () => dispatch("action", { type: "open_sqlite" }),
      },
      {
        id: "open_project",
        title: (await t("welcome.open_project")) || "開啟專案檔案",
        description:
          (await t("welcome.open_project_desc")) || "開啟 .insa 專案檔案",
        icon: "📁",
        action: () => dispatch("action", { type: "open_project" }),
      },
      {
        id: "new_project",
        title: (await t("welcome.new_project")) || "建立空白專案",
        description:
          (await t("welcome.new_project_desc")) || "開始一個新的空白專案",
        icon: "📄",
        action: () => {
          dispatch("action", { type: "new_project" });
        },
      },
    ];
  }
  // 組件初始化
  import { onMount } from "svelte";
  onMount(async () => {
    await initOptions();
    // 強制重新渲染
    options = [...options];
  });
</script>

<div class="welcome-container">
  <!-- 左側品牌區域 -->
  <div class="brand-section">
    <div class="logo-area">
      <img
        src={logoUrl}
        alt="Insyra Insights"
        class="logo-image"
        style="width: 100%;"
      />
      <h1 class="app-title">{appTitle}</h1>
      <p class="app-subtitle">{appSubtitle}</p>
    </div>
  </div>

  <!-- 右側操作區域 -->
  <div class="content-section">
    <div class="section-header">
      <h2>開始使用</h2>
    </div>

    <!-- 建立新專案 -->
    <div class="action-section">
      <h3>建立新專案</h3>
      <button
        class="action-card primary-action"
        on:click={() => dispatch("action", { type: "new_project" })}
      >
        <div class="action-icon">📄</div>
        <div class="action-content">
          <h4>新增空白專案</h4>
          <p>開始一個新的資料分析專案</p>
        </div>
      </button>
    </div>

    <!-- 開啟資料 -->
    <div class="action-section">
      <h3>開啟資料</h3>
      <div class="actions-grid">
        {#each options.slice(0, 4) as option}
          <button
            class="action-card"
            on:click={option.action}
            title={option.description}
          >
            <div class="action-icon">{option.icon}</div>
            <div class="action-content">
              <h4>{option.title}</h4>
              <p>{option.description}</p>
            </div>
          </button>
        {/each}
      </div>
    </div>

    <!-- 最近使用 -->
    <div class="action-section">
      <h3>{recentFilesTitle}</h3>
      <div class="recent-area">
        <div class="no-recent">{noRecentFiles}</div>
      </div>
    </div>
  </div>
</div>

<style>
  .welcome-container {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: linear-gradient(
      135deg,
      #87ceeb 0%,
      /* 較淺的天藍色 */ #7ac5cd 50%,
      /* 中間的天藍色 */ #6ca6cd 100% /* 較深的天藍色 */
    );
    display: grid;
    grid-template-columns: 1fr 2fr;
    z-index: 1000;
    font-family: /* 移除 "Nunito" */ -apple-system, BlinkMacSystemFont,
      "Segoe UI", Roboto, sans-serif;
    overflow: hidden;
    box-sizing: border-box;
  }
  /* 左側品牌區域 */
  .brand-section {
    background: linear-gradient(
      135deg,
      #6fb5d2 0%,
      /* 品牌區較亮的天藍色 */ #5a9ecb 100% /* 品牌區較深的天藍色 */
    );
    color: white; /* 文字顏色在深色背景上保持白色以確保對比度 */
    padding: 2rem 1.5rem; /* 減少 padding */
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    text-align: center;
    position: relative;
    overflow: hidden;
  }

  .brand-section::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-image: radial-gradient(
        circle at 20% 20%,
        rgba(255, 255, 255, 0.1) 1px,
        transparent 1px
      ),
      radial-gradient(
        circle at 80% 80%,
        rgba(255, 255, 255, 0.1) 1px,
        transparent 1px
      ),
      radial-gradient(
        circle at 40% 60%,
        rgba(255, 255, 255, 0.05) 1px,
        transparent 1px
      );
    background-size:
      100px 100px,
      80px 80px,
      120px 120px;
    animation: float 20s ease-in-out infinite;
    pointer-events: none;
  }

  .logo-area {
    position: relative;
    z-index: 1;
  }

  .logo-image {
    width: 120px;
    height: 120px;
    object-fit: contain;
    margin-bottom: 1.5rem; /* 減少 margin-bottom */
    filter: none; /* 移除logo陰影 */
  }

  .app-title {
    font-size: 3.5rem;
    font-weight: 700;
    color: white; /* 文字顏色在深色背景上保持白色 */
    margin: 0 0 0.75rem 0; /* 減少 margin-bottom */
    text-shadow: none; /* 移除文字陰影 */
  }

  .app-subtitle {
    font-size: 1.3rem;
    color: rgba(255, 255, 255, 0.85); /* 調整副標題顏色以搭配新的背景 */
    margin: 0;
    font-weight: 300;
  }

  /* 右側內容區域 */
  .content-section {
    background: #ffffff;
    padding: 3rem;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 2.5rem;
  }

  .section-header {
    border-bottom: 2px solid #e5e7eb;
    padding-bottom: 1rem;
  }

  .section-header h2 {
    font-size: 2.5rem;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }

  .action-section {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .action-section h3 {
    font-size: 1.5rem;
    font-weight: 600;
    color: #374151;
    margin: 0;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid #e5e7eb;
  }

  .action-card {
    background: #ffffff;
    border: 2px solid #e5e7eb;
    border-radius: 12px;
    padding: 1.5rem;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    display: flex;
    align-items: center;
    text-align: left;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    position: relative;
    overflow: hidden;
    font-family: inherit;
    width: 100%;
  }

  .action-card.primary-action {
    background: linear-gradient(
      135deg,
      #7ac5cd 0%,
      /* 主要按鈕較淺的天藍色 */ #6ca6cd 100% /* 主要按鈕較深的天藍色 */
    );
    color: white;
    border-color: #6ca6cd; /* 主要按鈕邊框 - 天藍色 */
    box-shadow: 0 4px 12px rgba(108, 166, 205, 0.3); /* 主要按鈕陰影 - 天藍色基底 */
  }

  .action-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
    border-color: #6ca6cd; /* 滑鼠懸停時的天藍色邊框 */
  }

  .action-card.primary-action:hover {
    box-shadow: 0 8px 25px rgba(108, 166, 205, 0.4); /* 主要按鈕滑鼠懸停陰影 - 天藍色基底 */
  }

  .actions-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1.5rem;
  }

  .action-icon {
    font-size: 2.5rem;
    margin-right: 1.5rem;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 60px;
    height: 60px;
    background: linear-gradient(135deg, #f3f4f6 0%, #e5e7eb 100%);
    border-radius: 12px;
  }

  .action-card.primary-action .action-icon {
    background: rgba(255, 255, 255, 0.2);
  }

  .action-content {
    flex: 1;
  }

  .action-content h4 {
    font-size: 1.2rem;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 0.5rem 0;
  }

  .action-card.primary-action .action-content h4 {
    color: white;
  }

  .action-content p {
    font-size: 0.95rem;
    color: #6b7280;
    margin: 0;
    line-height: 1.5;
  }

  .action-card.primary-action .action-content p {
    color: rgba(255, 255, 255, 0.8);
  }

  .recent-area {
    background: #f9fafb;
    border: 2px dashed #d1d5db;
    border-radius: 12px;
    padding: 3rem;
    text-align: center;
  }

  .no-recent {
    color: #9ca3af;
    font-style: italic;
    font-size: 1.1rem;
  }
  @keyframes float {
    0%,
    100% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(-10px);
    }
  }

  /* 響應式設計 */
  @media (max-width: 1200px) {
    .welcome-container {
      grid-template-columns: 1fr 1.8fr;
    }

    .content-section {
      padding: 2rem;
    }

    .brand-section {
      padding: 2rem 1.5rem;
    }

    .app-title {
      font-size: 3rem;
    }
  }

  @media (max-width: 768px) {
    .welcome-container {
      grid-template-columns: 1fr;
      grid-template-rows: auto 1fr;
    }

    .brand-section {
      padding: 2rem 1rem;
    }

    .app-title {
      font-size: 2.5rem;
    }

    .app-subtitle {
      font-size: 1.1rem;
    }

    .logo-image {
      width: 100px;
      height: 100px;
      margin-bottom: 1.5rem;
    }

    .content-section {
      padding: 2rem 1rem;
      gap: 2rem;
    }

    .section-header h2 {
      font-size: 2rem;
    }

    .actions-grid {
      grid-template-columns: 1fr;
      gap: 1rem;
    }

    .action-card {
      padding: 1.2rem;
    }

    .action-icon {
      font-size: 2rem;
      width: 50px;
      height: 50px;
      margin-right: 1rem;
    }
  }

  @media (max-width: 480px) {
    .brand-section {
      padding: 1.5rem 1rem;
    }

    .app-title {
      font-size: 2rem;
    }

    .app-subtitle {
      font-size: 1rem;
    }

    .logo-image {
      width: 80px;
      height: 80px;
      margin-bottom: 1rem;
    }

    .content-section {
      padding: 1.5rem 1rem;
      gap: 1.5rem;
    }

    .section-header h2 {
      font-size: 1.8rem;
    }

    .action-section h3 {
      font-size: 1.3rem;
    }

    .action-card {
      padding: 1rem;
      flex-direction: column;
      text-align: center;
    }

    .action-icon {
      margin-right: 0;
      margin-bottom: 1rem;
    }

    .action-content h4 {
      font-size: 1.1rem;
    }

    .action-content p {
      font-size: 0.9rem;
    }
  }
</style>
