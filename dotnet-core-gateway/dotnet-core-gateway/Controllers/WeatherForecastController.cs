using Microsoft.AspNetCore.Mvc;

namespace dotnet_core_gateway.Controllers
{
    [ApiController]
    [Route("[controller]")]
    public class WeatherForecastController : ControllerBase
    {
        private static readonly string[] Summaries = new[]
        {
            "Freezing", "Bracing", "Chilly", "Cool", "Mild", "Warm", "Balmy", "Hot", "Sweltering", "Scorching"
        };

        private readonly ILogger<WeatherForecastController> _logger;

        private readonly HttpClient _httpClient;

        public WeatherForecastController(ILogger<WeatherForecastController> logger, HttpClient httpClient)
        {
            _logger = logger;
            _httpClient = httpClient;
        }

        [HttpGet(Name = "GetWeatherForecast")]
        public async Task<ApiResponse<WeatherForecast[]>> Get()
        {
            await Task.Delay(500);
            var weathers = await _httpClient.GetFromJsonAsync<ApiResponse<WeatherForecast[]>>("http://localhost:8080/weather");
            if (weathers == null) throw new InvalidOperationException("Failed to get weather");
            await Task.Delay(800);
            return weathers;
        }
    }
}