<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { GetText } from "../../wailsjs/go/main/App";

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
  let appTitle = "Insyra Insights";
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
    appTitle = (await t("welcome.title")) || "Insyra Insights";
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
  <div class="welcome-content">
    <!-- 應用程式標題和Logo -->
    <div class="header">
      <div class="logo">
        <img
          src="/src/assets/images/logo-universal.png"
          alt="Insyra Insights"
          class="logo-image"
        />
      </div>
      <h1 class="app-title">{appTitle}</h1>
      <p class="app-subtitle">{appSubtitle}</p>
    </div>

    <!-- 歡迎選項 -->
    <div class="options-grid">
      {#each options as option}
        <button
          class="option-card"
          on:click={option.action}
          title={option.description}
        >
          <div class="option-icon">{option.icon}</div>
          <div class="option-content">
            <h3 class="option-title">{option.title}</h3>
            <p class="option-description">{option.description}</p>
          </div>
        </button>
      {/each}
    </div>
    <!-- 最近使用的檔案 -->
    <div class="recent-files">
      <h3>{recentFilesTitle}</h3>
      <div class="recent-list">
        <!-- 這裡可以添加最近使用的檔案列表 -->
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
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    font-family:
      "Nunito",
      -apple-system,
      BlinkMacSystemFont,
      "Segoe UI",
      Roboto,
      sans-serif;
  }

  .welcome-content {
    max-width: 900px;
    width: 90%;
    background: rgba(255, 255, 255, 0.95);
    border-radius: 20px;
    padding: 3rem;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(10px);
    animation: fadeInUp 0.6s ease-out;
  }

  @keyframes fadeInUp {
    from {
      opacity: 0;
      transform: translateY(30px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .header {
    text-align: center;
    margin-bottom: 3rem;
  }

  .logo {
    margin-bottom: 1rem;
  }

  .logo-image {
    width: 80px;
    height: 80px;
    object-fit: contain;
  }

  .app-title {
    font-size: 2.5rem;
    font-weight: 700;
    color: #2d3748;
    margin: 0 0 0.5rem 0;
  }

  .app-subtitle {
    font-size: 1.1rem;
    color: #718096;
    margin: 0;
  }

  .options-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1.5rem;
    margin-bottom: 3rem;
  }

  .option-card {
    background: white;
    border: 2px solid #e2e8f0;
    border-radius: 12px;
    padding: 1.5rem;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    text-align: left;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  }

  .option-card:hover {
    border-color: #667eea;
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(102, 126, 234, 0.15);
  }

  .option-card:active {
    transform: translateY(0);
  }

  .option-icon {
    font-size: 2.5rem;
    margin-right: 1.2rem;
    flex-shrink: 0;
  }

  .option-content {
    flex: 1;
  }

  .option-title {
    font-size: 1.2rem;
    font-weight: 600;
    color: #2d3748;
    margin: 0 0 0.5rem 0;
  }

  .option-description {
    font-size: 0.9rem;
    color: #718096;
    margin: 0;
    line-height: 1.4;
  }

  .recent-files {
    border-top: 1px solid #e2e8f0;
    padding-top: 2rem;
  }

  .recent-files h3 {
    font-size: 1.3rem;
    font-weight: 600;
    color: #2d3748;
    margin: 0 0 1rem 0;
  }

  .recent-list {
    min-height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .no-recent {
    color: #a0aec0;
    font-style: italic;
  }

  /* 響應式設計 */
  @media (max-width: 768px) {
    .welcome-content {
      padding: 2rem;
      margin: 1rem;
    }

    .app-title {
      font-size: 2rem;
    }

    .options-grid {
      grid-template-columns: 1fr;
      gap: 1rem;
    }

    .option-card {
      padding: 1.2rem;
    }

    .option-icon {
      font-size: 2rem;
      margin-right: 1rem;
    }
  }
</style>
