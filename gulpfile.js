'use strict';

const browserify = require('browserify')
const gulp = require('gulp')
const source = require('vinyl-source-stream')
const buffer = require('vinyl-buffer')
const rename = require('gulp-rename')
const gutil = require('gulp-util')
const sass = require('gulp-sass')
const uglify = require('gulp-uglify')
const babel = require('gulp-babel');
const cachebust = require('gulp-cache-bust');

const debug = process.env.NODE_ENV !== 'production';
let defaultTasks = [ 'app-js', 'tracker-js', 'sass', 'html', 'img', 'fonts' ] ;

gulp.task('default', defaultTasks);

gulp.task('app-js', function () {
    let stream = browserify({
        entries: './assets/src/js/script.js',
        debug: debug
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
  
    if(!debug) {
      stream.pipe(buffer()).pipe(uglify())
    }
    
    return stream.pipe(gulp.dest(`./assets/build/js`))  
});  

gulp.task('tracker-js', function () {
  return gulp.src('./assets/src/js/tracker.js')
        .pipe(babel({
            presets: ["@babel/preset-env"],
        }))
        .pipe(uglify())
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

gulp.task('sass', function () {
	var files = './assets/src/sass/[^_]*.scss';
	return gulp.src(files)
		.pipe(sass())
    .on('error', gutil.log)
		.pipe(rename({ extname: '.css' }))
		.pipe(gulp.dest(`./assets/build/css`))
});

gulp.task('watch', ['default'], function() {
  gulp.watch(['./assets/src/js/**/*.js'], ['app-js', 'tracker-js'] );
  gulp.watch(['./assets/src/sass/**/**/*.scss'], ['sass'] );
  gulp.watch(['./assets/src/**/*.html'], ['html'] );
  gulp.watch(['./assets/src/img/**/*'], ['img'] );
  gulp.watch(['./assets/src/fonts/**/*'], ['fonts'] );
});
