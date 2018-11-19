'use strict';

const browserify = require('browserify')
const gulp = require('gulp')
const source = require('vinyl-source-stream')
const buffer = require('vinyl-buffer')
const uglify = require('gulp-uglify')
const babel = require('gulp-babel');
const cachebust = require('gulp-cache-bust');
const concat = require('gulp-concat');
const gulpif = require('gulp-if')

const debug = process.env.NODE_ENV !== 'production';

gulp.task('app-js', function () {
    return browserify({
        entries: './assets/src/js/script.js',
        debug: debug,
        ignoreMissing: true,
    })
    .transform("babelify", {
      presets: ["@babel/preset-env"],
      plugins: [ 
        ["@babel/plugin-proposal-decorators", { "legacy": true }],
        ["@babel/plugin-transform-react-jsx", { "pragma":"h" } ]
      ]
    })    
    .bundle()  
    .pipe(source('script.js'))
	.pipe(buffer())
	.pipe(gulpif(!debug, uglify()))
  	.pipe(gulp.dest(`./assets/build/js`))  
});  

gulp.task('tracker-js', function () {
  return gulp.src('./assets/src/js/tracker.js')
        .pipe(babel({
            presets: ["@babel/preset-env"],
        }))
        .pipe(gulpif(!debug, uglify()))
        .pipe(gulp.dest('./assets/build/js'));
});

gulp.task('fonts', function() {
  return gulp.src('./assets/src/fonts/**/*')
    .pipe(gulp.dest(`./assets/build/fonts`))
});

gulp.task('img', function() {
  return gulp.src('./assets/src/img/**/*')
    .pipe(gulp.dest(`./assets/build/img`))
});

gulp.task('html', function() {
  return gulp.src('./assets/src/**/*.html')
    .pipe(cachebust({
      type: 'timestamp'
    }))
    .pipe(gulp.dest(`./assets/build/`))
});

gulp.task('css', function () {
	return gulp.src('./assets/src/css/*.css')
		.pipe(concat('styles.css'))
		.pipe(gulp.dest(`./assets/build/css`))
});

gulp.task('default', gulp.series('app-js', 'tracker-js', 'css', 'html', 'img', 'fonts' ) );

gulp.task('watch', gulp.series('default', function() {
  gulp.watch(['./assets/src/js/**/*.js'], gulp.parallel('app-js', 'tracker-js') );
  gulp.watch(['./assets/src/sass/**/**/*.scss'], gulp.parallel( 'css') );
  gulp.watch(['./assets/src/**/*.html'], gulp.parallel( 'html') );
  gulp.watch(['./assets/src/img/**/*'], gulp.parallel( 'img') );
  gulp.watch(['./assets/src/fonts/**/*'], gulp.parallel( 'fonts') );
}));
